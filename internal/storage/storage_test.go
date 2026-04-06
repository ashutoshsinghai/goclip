package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// withTempHome redirects HistoryFile to a temp directory for the duration of t.
func withTempHome(t *testing.T) {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp) // Windows
}

// makeClip is a convenience constructor for test clips.
func makeClip(id int, content string, pinned bool) Clip {
	return Clip{ID: id, Content: content, CopiedAt: time.Now(), Pinned: pinned}
}

// -----------------------------------------------------------------------
// AddClip
// -----------------------------------------------------------------------

func TestAddClip_FirstEntry(t *testing.T) {
	clips := AddClip("hello", []Clip{})
	if len(clips) != 1 {
		t.Fatalf("expected 1 clip, got %d", len(clips))
	}
	if clips[0].Content != "hello" {
		t.Errorf("expected content %q, got %q", "hello", clips[0].Content)
	}
	if clips[0].ID != 1 {
		t.Errorf("expected ID 1, got %d", clips[0].ID)
	}
}

func TestAddClip_PrependsBefore_Existing(t *testing.T) {
	existing := []Clip{makeClip(1, "first", false)}
	clips := AddClip("second", existing)
	if clips[0].Content != "second" {
		t.Errorf("expected newest entry at index 0, got %q", clips[0].Content)
	}
	if clips[1].Content != "first" {
		t.Errorf("expected old entry at index 1, got %q", clips[1].Content)
	}
}

func TestAddClip_IDIncrements(t *testing.T) {
	clips := AddClip("a", []Clip{})
	clips = AddClip("b", clips)
	clips = AddClip("c", clips)
	if clips[0].ID != 3 {
		t.Errorf("expected ID 3, got %d", clips[0].ID)
	}
}

func TestAddClip_DuplicateSkipped(t *testing.T) {
	clips := AddClip("hello", []Clip{})
	clips = AddClip("hello", clips)
	if len(clips) != 1 {
		t.Errorf("expected duplicate to be skipped, got %d clips", len(clips))
	}
}

func TestAddClip_DuplicateNotSkipped_WhenNotMostRecent(t *testing.T) {
	clips := AddClip("hello", []Clip{})
	clips = AddClip("world", clips)
	clips = AddClip("hello", clips) // same as index 1, but not index 0
	if len(clips) != 3 {
		t.Errorf("expected 3 clips, got %d", len(clips))
	}
}

func TestAddClip_CapsAtMaxHistory(t *testing.T) {
	clips := []Clip{}
	for i := 0; i < MaxHistory+10; i++ {
		clips = AddClip(string(rune('a'+i%26))+string(rune(i)), clips)
	}
	unpinned := 0
	for _, c := range clips {
		if !c.Pinned {
			unpinned++
		}
	}
	if unpinned > MaxHistory {
		t.Errorf("expected at most %d unpinned clips, got %d", MaxHistory, unpinned)
	}
}

func TestAddClip_PinnedItemsNotEvicted(t *testing.T) {
	clips := []Clip{}
	for i := 0; i < MaxHistory; i++ {
		clips = AddClip(string(rune('a'+i%26))+string(rune(i)), clips)
	}
	oldest := clips[len(clips)-1]
	clips, _ = TogglePin(oldest.ID, clips)

	for i := 0; i < 5; i++ {
		clips = AddClip("overflow"+string(rune(i)), clips)
	}

	found := false
	for _, c := range clips {
		if c.ID == oldest.ID && c.Pinned {
			found = true
			break
		}
	}
	if !found {
		t.Error("pinned clip was incorrectly evicted")
	}
}

// -----------------------------------------------------------------------
// TogglePin
// -----------------------------------------------------------------------

func TestTogglePin_PinsUnpinnedClip(t *testing.T) {
	clips := []Clip{makeClip(1, "test", false)}
	clips, pinned := TogglePin(1, clips)
	if !pinned {
		t.Error("expected TogglePin to return true (now pinned)")
	}
	if !clips[0].Pinned {
		t.Error("expected clip to be pinned")
	}
}

func TestTogglePin_UnpinsPinnedClip(t *testing.T) {
	clips := []Clip{makeClip(1, "test", true)}
	clips, pinned := TogglePin(1, clips)
	if pinned {
		t.Error("expected TogglePin to return false (now unpinned)")
	}
	if clips[0].Pinned {
		t.Error("expected clip to be unpinned")
	}
}

func TestTogglePin_MissingID(t *testing.T) {
	clips := []Clip{makeClip(1, "test", false)}
	clips, pinned := TogglePin(99, clips)
	if pinned {
		t.Error("expected false for missing ID")
	}
	if clips[0].Pinned {
		t.Error("existing clip should not be modified")
	}
}

func TestTogglePin_OnlyAffectsTargetID(t *testing.T) {
	clips := []Clip{
		makeClip(1, "one", false),
		makeClip(2, "two", false),
		makeClip(3, "three", false),
	}
	clips, _ = TogglePin(2, clips)
	if clips[0].Pinned || clips[2].Pinned {
		t.Error("only clip #2 should be pinned")
	}
	if !clips[1].Pinned {
		t.Error("clip #2 should be pinned")
	}
}

// -----------------------------------------------------------------------
// Sorted
// -----------------------------------------------------------------------

func TestSorted_PinnedFirst(t *testing.T) {
	clips := []Clip{
		makeClip(1, "unpinned-a", false),
		makeClip(2, "pinned-a", true),
		makeClip(3, "unpinned-b", false),
		makeClip(4, "pinned-b", true),
	}
	sorted := Sorted(clips)

	for i := 0; i < 2; i++ {
		if !sorted[i].Pinned {
			t.Errorf("index %d should be pinned, got %q", i, sorted[i].Content)
		}
	}
	for i := 2; i < 4; i++ {
		if sorted[i].Pinned {
			t.Errorf("index %d should be unpinned, got %q", i, sorted[i].Content)
		}
	}
}

func TestSorted_PreservesRelativeOrder(t *testing.T) {
	clips := []Clip{
		makeClip(1, "unpinned-a", false),
		makeClip(2, "pinned-a", true),
		makeClip(3, "unpinned-b", false),
		makeClip(4, "pinned-b", true),
	}
	sorted := Sorted(clips)

	if sorted[0].ID != 2 || sorted[1].ID != 4 {
		t.Errorf("pinned relative order wrong: got IDs %d, %d", sorted[0].ID, sorted[1].ID)
	}
	if sorted[2].ID != 1 || sorted[3].ID != 3 {
		t.Errorf("unpinned relative order wrong: got IDs %d, %d", sorted[2].ID, sorted[3].ID)
	}
}

func TestSorted_EmptySlice(t *testing.T) {
	result := Sorted([]Clip{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d items", len(result))
	}
}

func TestSorted_AllPinned(t *testing.T) {
	clips := []Clip{
		makeClip(1, "a", true),
		makeClip(2, "b", true),
	}
	sorted := Sorted(clips)
	if sorted[0].ID != 1 || sorted[1].ID != 2 {
		t.Error("order should be preserved when all are pinned")
	}
}

// -----------------------------------------------------------------------
// Load / Save round-trip
// -----------------------------------------------------------------------

func TestSaveLoad_RoundTrip(t *testing.T) {
	withTempHome(t)

	original := []Clip{
		makeClip(2, "world", true),
		makeClip(1, "hello", false),
	}
	Save(original)

	loaded := Load()
	if len(loaded) != len(original) {
		t.Fatalf("expected %d clips, got %d", len(original), len(loaded))
	}
	for i, c := range original {
		if c.ID != loaded[i].ID || c.Content != loaded[i].Content || c.Pinned != loaded[i].Pinned {
			t.Errorf("clip %d mismatch: want %+v, got %+v", i, c, loaded[i])
		}
	}
}

func TestLoad_MissingFile(t *testing.T) {
	withTempHome(t)
	clips := Load()
	if clips == nil {
		t.Error("expected non-nil slice when file is missing")
	}
	if len(clips) != 0 {
		t.Errorf("expected empty slice, got %d clips", len(clips))
	}
}

func TestLoad_CorruptedFile(t *testing.T) {
	withTempHome(t)
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".goclip")
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "history.json"), []byte("not json {{{"), 0644)

	clips := Load()
	if clips == nil || len(clips) != 0 {
		t.Error("expected empty slice for corrupted file")
	}
}

func TestSave_EmptySlice(t *testing.T) {
	withTempHome(t)
	Save([]Clip{})
	clips := Load()
	if len(clips) != 0 {
		t.Errorf("expected 0 clips after saving empty slice, got %d", len(clips))
	}
}
