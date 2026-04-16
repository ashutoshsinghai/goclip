# goclip

A fast, terminal-native clipboard history manager written in Go.

Run a background daemon to capture everything you copy, then recall any entry instantly through a TUI picker or plain CLI commands.

## Features

- **Daemon** — polls the clipboard every 500ms and persists entries to `~/.goclip/history.json`
- **Menu bar / system tray** — native macOS menu bar app with history grouped by date; shows timestamps, fires a notification on copy, updates instantly via file-watch
- **Interactive TUI** — scrollable picker with live search and pin support, built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Pin items** — keep important clips at the top of your history
- **Non-interactive CLI** — `list`, `copy <id>`, `pin <id>`, `clear` for scripting
- **Keeps up to 100 entries**, de-duplicating consecutive identical clips

## Installation

### macOS / Linux

```bash
curl -sSL https://raw.githubusercontent.com/ashutoshsinghai/goclip/main/scripts/install.sh | sh
```

Installs to `/usr/local/bin` (or `~/bin` if no write access). Works on Intel and Apple Silicon.

---

### Windows — PowerShell

```powershell
irm https://raw.githubusercontent.com/ashutoshsinghai/goclip/main/scripts/install.ps1 | iex
```

Installs to `%USERPROFILE%\bin` and adds it to your PATH automatically. No admin needed.

---

### Windows — Command Prompt (cmd.exe)

Windows 10 and 11 have `curl` and `tar` built in. Run these commands:

```cmd
curl -L -o goclip.zip https://github.com/ashutoshsinghai/goclip/releases/latest/download/goclip_windows_amd64.zip
tar -xf goclip.zip
mkdir %USERPROFILE%\bin
move goclip.exe %USERPROFILE%\bin\
del goclip.zip
```

Then add `%USERPROFILE%\bin` to your PATH:

```cmd
setx PATH "%PATH%;%USERPROFILE%\bin"
```

Restart your terminal and run `goclip help` to verify.

---

### Manual download (any platform)

Download the right binary for your OS from [GitHub Releases](https://github.com/ashutoshsinghai/goclip/releases), extract it, and move it somewhere on your PATH.

| OS | File |
|---|---|
| macOS Apple Silicon | `goclip_darwin_arm64.tar.gz` |
| macOS Intel | `goclip_darwin_amd64.tar.gz` |
| Linux ARM64 | `goclip_linux_arm64.tar.gz` |
| Linux x86-64 | `goclip_linux_amd64.tar.gz` |
| Windows x86-64 | `goclip_windows_amd64.zip` |

---

### Build from source (requires Go 1.21+)

```bash
go install github.com/ashutoshsinghai/goclip@latest
```

## Usage

### Typical workflow

```bash
# 1. Start the daemon in the background
goclip daemon

# 2a. Open the menu bar app — lives in your menu bar, always one click away
goclip tray

# 2b. Or open the interactive TUI picker in the terminal
goclip pick
```

### All commands

| Command | Description |
|---|---|
| `goclip daemon` | Start clipboard watcher in background |
| `goclip stop` | Stop background daemon |
| `goclip status` | Show whether daemon is running |
| `goclip run` | Run clipboard watcher in foreground |
| `goclip tray` | Start menu bar / system tray app in background |
| `goclip tray stop` | Stop the tray app |
| `goclip tray status` | Show whether the tray app is running |
| `goclip pick` | Open the interactive TUI picker |
| `goclip list` | Print history as a plain text table |
| `goclip search <keyword>` | Search history by keyword |
| `goclip copy <id>` | Put a historical entry back on the clipboard |
| `goclip pin <id>` | Pin/unpin an entry (pinned items stay at the top) |
| `goclip clear` | Wipe all saved history |
| `goclip upgrade` | Upgrade to the latest version |
| `goclip uninstall` | Remove goclip from your system |
| `goclip version` | Show current version |
| `goclip help` | Show usage |

### Menu bar app (macOS & Windows)

`goclip tray` puts a clipboard icon in your menu bar. History is grouped by date — hover over a group to see its entries, click any entry to copy it instantly.

```
📌 Pinned  (2)    →  3:04 PM  ·  My SSH key
                     2:51 PM  ·  SELECT * FROM users…
Today  (8)        →  4:12 PM  ·  Hello world
                     3:58 PM  ·  npm install --save-dev
                     … and 6 more
Yesterday  (3)    →  Apr 2  11:03 AM  ·  git commit -m "fix:…
────────────────────────────────
🔍  Open Picker…          ← full TUI search in a terminal window
────────────────────────────────
Quit goclip tray
```

- Up to **20 entries per date group**; overflow shows `… and N more` which opens the full picker
- **macOS notification** fired on every copy
- Menu updates **instantly** when the daemon captures a new clip (file-watch, no polling delay)
- Survives terminal close; stop it with `goclip tray stop`

### TUI keybindings

| Key | Action |
|---|---|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Copy selected entry and exit |
| `p` | Pin / unpin selected entry |
| `/` | Enter search mode |
| `Esc` | Exit search mode |
| `q` | Quit |

## Upgrading

```bash
goclip upgrade
```

Detects your OS and architecture, downloads the latest binary from GitHub Releases, and replaces the current binary in-place. No need to re-run the install script.

---

## Uninstalling

### macOS / Linux

```bash
rm $(which goclip)       # remove the binary
rm -rf ~/.goclip         # remove saved history (optional)
```

### Windows — PowerShell

```powershell
Remove-Item "$env:USERPROFILE\bin\goclip.exe"
Remove-Item -Recurse "$env:USERPROFILE\.goclip"   # optional: remove history
```

### Windows — Command Prompt

```cmd
del %USERPROFILE%\bin\goclip.exe
rmdir /s %USERPROFILE%\.goclip
```

---

## Reinstalling

Just run the install command for your platform again — it will overwrite the existing binary.

---

## Storage

History is stored at `~/.goclip/history.json`. Plain JSON — you can inspect, back up, or edit it manually.

## Running the daemon

### Quick start (background)

```bash
goclip daemon     # starts in background, detached from terminal
goclip status     # check it's running
goclip stop       # stop it
```

Logs are written to `~/.goclip/daemon.log`:
```bash
tail -f ~/.goclip/daemon.log
```

### Foreground mode (for debugging)

```bash
goclip run        # runs in terminal, Ctrl+C to stop
```

### Auto-start on login

#### macOS (launchd)

Create `~/Library/LaunchAgents/com.goclip.daemon.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.goclip.daemon</string>
  <key>ProgramArguments</key>
  <array>
    <string>/usr/local/bin/goclip</string>
    <string>run</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <true/>
</dict>
</plist>
```

```bash
launchctl load ~/Library/LaunchAgents/com.goclip.daemon.plist
```

#### Linux (systemd)

Create `~/.config/systemd/user/goclip.service`:

```ini
[Unit]
Description=goclip clipboard daemon

[Service]
ExecStart=/usr/local/bin/goclip run
Restart=on-failure

[Install]
WantedBy=default.target
```

```bash
systemctl --user enable --now goclip
```

## Project structure

```
goclip/
├── main.go                    # CLI entry point, argument routing
├── cmd/
│   ├── daemon/
│   │   ├── daemon.go          # run, daemon, stop, status
│   │   ├── daemon_unix.go     # background process (macOS/Linux)
│   │   └── daemon_windows.go  # background process (Windows)
│   ├── tray/
│   │   ├── tray.go            # menu bar UI, start/stop/status, date grouping
│   │   ├── tray_unix.go       # spawn/kill tray process (macOS/Linux)
│   │   ├── tray_windows.go    # spawn/kill tray process (Windows)
│   │   ├── notify_darwin.go   # macOS notification via osascript + open picker
│   │   ├── notify_windows.go  # Windows toast notification + open picker
│   │   ├── notify_other.go    # no-op stubs for other platforms
│   │   ├── signal_unix.go     # ignore SIGHUP so tray survives terminal close
│   │   └── signal_windows.go  # no-op on Windows
│   ├── history.go             # list, search
│   ├── clip.go                # copy, pin, clear
│   ├── upgrade.go             # self-upgrade from GitHub Releases
│   ├── extract.go             # tar.gz / zip extraction helpers
│   └── uninstall.go           # self-uninstall
├── internal/
│   ├── storage/
│   │   └── storage.go         # read/write history.json
│   └── style/
│       └── style.go           # shared terminal styles
├── ui/
│   └── tui.go                 # Bubbletea TUI picker
└── scripts/
    ├── install.sh             # installer for macOS/Linux
    └── install.ps1            # installer for Windows (PowerShell)
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — [Ashutosh Singhai](https://github.com/ashutoshsinghai)
