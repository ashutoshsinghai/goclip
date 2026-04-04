# goclip

A fast, terminal-native clipboard history manager written in Go.

Run a background daemon to capture everything you copy, then recall any entry instantly through a TUI picker or plain CLI commands.

## Features

- **Daemon** — polls the clipboard every 500ms and persists entries to `~/.goclip/history.json`
- **Interactive TUI** — scrollable picker with live search and pin support, built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Pin items** — keep important clips at the top of your history
- **Non-interactive CLI** — `list`, `copy <id>`, `pin <id>`, `clear` for scripting
- **Keeps up to 100 entries**, de-duplicating consecutive identical clips

## Installation

### Option 1 — Install script (no Go needed)

```bash
curl -sSL https://raw.githubusercontent.com/ashutoshsinghai/goclip/main/install.sh | sh
```

### Option 2 — Download binary manually

Grab the latest binary for your platform from [GitHub Releases](https://github.com/ashutoshsinghai/goclip/releases), extract it, and move it to your PATH.

### Option 3 — Build from source (requires Go 1.21+)

```bash
go install github.com/ashutoshsinghai/goclip@latest
```

## Usage

### Typical workflow

```bash
# 1. Start the daemon in a background terminal tab
goclip daemon

# 2. Open the interactive picker whenever you need an old clip
goclip pick
```

### All commands

| Command | Description |
|---|---|
| `goclip daemon` | Start watching the clipboard |
| `goclip pick` | Open the interactive TUI picker |
| `goclip list` | Print history as a plain text table |
| `goclip copy <id>` | Put a historical entry back on the clipboard |
| `goclip pin <id>` | Pin/unpin an entry (pinned items stay at the top) |
| `goclip clear` | Wipe all saved history |
| `goclip help` | Show usage |

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

## Storage

History is stored at `~/.goclip/history.json`. Plain JSON — you can inspect, back up, or edit it manually.

## Running as a background service

### macOS (launchd)

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
    <string>daemon</string>
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

### Linux (systemd)

Create `~/.config/systemd/user/goclip.service`:

```ini
[Unit]
Description=goclip clipboard daemon

[Service]
ExecStart=/usr/local/bin/goclip daemon
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
├── main.go            # CLI entry point and argument routing
├── commands/
│   ├── daemon.go      # Clipboard polling loop
│   └── list.go        # list, copy, pin, clear subcommands
├── storage/
│   └── storage.go     # Read/write history.json
└── ui/
    └── tui.go         # Bubbletea TUI picker
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — [Ashutosh Singhai](https://github.com/ashutoshsinghai)
