//go:build darwin || windows

// Package tray implements a system tray / menu bar app for goclip.
//
// Menu structure:
//
//	📌 Pinned (N)   →  item…
//	                   item…
//	Today (N)       →  item…
//	Yesterday (N)   →  item…
//	This Week (N)   →  item…
//	Older           →  item…
//	────────────────
//	🔍 Open Picker…
//	────────────────
//	Quit
package tray

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"fyne.io/systray"
	"github.com/ashutoshsinghai/goclip/internal/storage"
	"github.com/fsnotify/fsnotify"
	"github.com/atotto/clipboard"
)

const (
	maxPinnedSlots         = 5
	maxSlotsPerGroup       = 20
	refreshInterval        = time.Second
	windowsBoundaryRefresh = 5 * time.Minute
)

func traySupported() bool { return true }

// ── Date grouping ─────────────────────────────────────────────────────────────

type bucketKey int

const (
	bucketPinned bucketKey = iota
	bucketToday
	bucketYesterday
	bucketThisWeek
	bucketOlder
)

func bucketFor(c storage.Clip) bucketKey {
	if c.Pinned {
		return bucketPinned
	}
	now := time.Now()
	start := func(d time.Time) time.Time {
		return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	}
	today := start(now)
	yesterday := today.Add(-24 * time.Hour)
	weekAgo := today.Add(-7 * 24 * time.Hour)

	t := c.CopiedAt
	switch {
	case !t.Before(today):
		return bucketToday
	case !t.Before(yesterday):
		return bucketYesterday
	case !t.Before(weekAgo):
		return bucketThisWeek
	default:
		return bucketOlder
	}
}

// ── Slot / group types ───────────────────────────────────────────────────────

type slot struct {
	item    *systray.MenuItem
	mu      sync.Mutex
	content string
	title   string  // full menu label: "3:04 PM  ·  Hello world…"
	preview string  // raw content only, used for notifications (no time prefix)
}

// group is one collapsible date bucket in the menu.
// Its header is a top-level item; slots are submenus of the header.
// overflow is an extra subitem shown when there are more clips than slots.
type group struct {
	header    *systray.MenuItem
	slots     []*slot
	overflow  *systray.MenuItem // "… and N more" — opens the full picker
	baseLabel string
	bucket    bucketKey
}

func newGroup(baseLabel string, size int, bucket bucketKey) *group {
	g := &group{
		baseLabel: baseLabel,
		slots:     make([]*slot, size),
		bucket:    bucket,
	}
	g.header = systray.AddMenuItem(baseLabel, "")
	g.header.Hide()
	for i := range g.slots {
		s := &slot{item: g.header.AddSubMenuItem("", "")}
		s.item.Hide()
		g.slots[i] = s
		go listenForClick(s)
	}
	// Overflow item — always last in the submenu
	g.overflow = g.header.AddSubMenuItem("", "")
	g.overflow.Hide()
	go func() {
		for range g.overflow.ClickedCh {
			openPicker()
		}
	}()
	return g
}

// fill updates the group with the given clips, showing/hiding as needed.
func (g *group) fill(clips []storage.Clip) {
	total := len(clips)
	if total == 0 {
		g.header.Hide()
		return
	}

	g.header.SetTitle(fmt.Sprintf("%s  (%d)", g.baseLabel, total))
	g.header.Show()

	shown := total
	if shown > len(g.slots) {
		shown = len(g.slots)
	}

	for i, s := range g.slots {
		if i < shown {
			c := clips[i]
			raw := contentPreview(c)
			title := timeLabel(c, g.bucket) + "  ·  " + raw
			s.mu.Lock()
			s.content = c.Content
			s.title = title
			s.preview = raw
			s.mu.Unlock()
			s.item.SetTitle(title)
			s.item.Show()
		} else {
			s.mu.Lock()
			s.content = ""
			s.title = ""
			s.preview = ""
			s.mu.Unlock()
			s.item.Hide()
		}
	}

	extra := total - shown
	if extra > 0 {
		g.overflow.SetTitle(fmt.Sprintf("  … and %d more", extra))
		g.overflow.Show()
	} else {
		g.overflow.Hide()
	}
}

var groups map[bucketKey]*group

// ── Tray UI ──────────────────────────────────────────────────────────────────

// Run starts the systray UI. Blocks until the user clicks Quit.
func Run() {
	ignoreSighup()
	systray.Run(onReady, nil)
}

func onReady() {
	icon := clipboardIcon()
	systray.SetTemplateIcon(icon, icon)
	systray.SetTooltip("goclip — Clipboard History")

	groups = map[bucketKey]*group{
		bucketPinned:    newGroup("📌  Pinned", maxPinnedSlots, bucketPinned),
		bucketToday:     newGroup("Today", maxSlotsPerGroup, bucketToday),
		bucketYesterday: newGroup("Yesterday", maxSlotsPerGroup, bucketYesterday),
		bucketThisWeek:  newGroup("This Week", maxSlotsPerGroup, bucketThisWeek),
		bucketOlder:     newGroup("Older", maxSlotsPerGroup, bucketOlder),
	}

	systray.AddSeparator()
	picker := systray.AddMenuItem("🔍  Open Picker...", "Search and browse all clipboard history")
	go func() {
		for range picker.ClickedCh {
			openPicker()
		}
	}()

	systray.AddSeparator()
	quit := systray.AddMenuItem("Quit goclip tray", "Stop the menu bar app")
	go func() {
		<-quit.ClickedCh
		systray.Quit()
	}()

	refresh()
	go watchHistory()
	go watchBoundary()
}

// watchBoundary re-buckets clips when the day rolls over.
// On macOS, fyne/systray emits TrayOpenedCh on every menu open — refreshing
// then is enough and costs nothing while idle. On Windows the library never
// signals menu opens, so we fall back to a slow ticker.
func watchBoundary() {
	go func() {
		for range systray.TrayOpenedCh {
			refresh()
		}
	}()

	if runtime.GOOS == "windows" {
		ticker := time.NewTicker(windowsBoundaryRefresh)
		defer ticker.Stop()
		for range ticker.C {
			refresh()
		}
	}
}

// watchHistory watches history.json and calls refresh whenever it changes.
// Falls back to polling every refreshInterval if the watcher can't be set up.
func watchHistory() {
	histFile := storage.HistoryFile()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		// fallback: poll
		for {
			time.Sleep(refreshInterval)
			refresh()
		}
	}
	defer watcher.Close()

	if err := watcher.Add(histFile); err != nil {
		// file might not exist yet — poll until it does, then switch to watching
		for {
			time.Sleep(refreshInterval)
			refresh()
			if watcher.Add(histFile) == nil {
				break
			}
		}
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				refresh()
			}
		case _, ok := <-watcher.Errors:
			if !ok {
				return
			}
		}
	}
}

func listenForClick(s *slot) {
	for range s.item.ClickedCh {
		s.mu.Lock()
		text := s.content
		preview := s.preview // raw content only — no time prefix in the notification
		s.mu.Unlock()
		if text == "" {
			continue
		}
		clipboard.WriteAll(text) //nolint:errcheck
		notifyCopied(preview)
	}
}

func refresh() {
	clips := storage.Sorted(storage.Load())

	buckets := map[bucketKey][]storage.Clip{
		bucketPinned:    nil,
		bucketToday:     nil,
		bucketYesterday: nil,
		bucketThisWeek:  nil,
		bucketOlder:     nil,
	}
	for _, c := range clips {
		k := bucketFor(c)
		buckets[k] = append(buckets[k], c)
	}

	// Fill in display order
	for _, k := range []bucketKey{bucketPinned, bucketToday, bucketYesterday, bucketThisWeek, bucketOlder} {
		groups[k].fill(buckets[k])
	}
}

// contentPreview returns a clean single-line preview of the clip (no timestamp).
func contentPreview(c storage.Clip) string {
	text := strings.ReplaceAll(c.Content, "\n", " ↵ ")
	text = strings.ReplaceAll(text, "\t", " → ")
	text = strings.TrimSpace(text)
	runes := []rune(text)
	if len(runes) > 45 {
		text = string(runes[:45]) + "…"
	}
	return text
}

// timeLabel returns the timestamp prefix for a clip based on which bucket it's in.
// Today shows time only; everything else shows date + time.
func timeLabel(c storage.Clip, k bucketKey) string {
	switch k {
	case bucketToday, bucketPinned:
		return c.CopiedAt.Format("3:04 PM")
	default:
		return c.CopiedAt.Format("Jan 2  3:04 PM")
	}
}

