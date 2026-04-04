# goclip

A fast, terminal-native clipboard history manager written in Go.

Run a background daemon to capture everything you copy, then recall any entry instantly through a fuzzy-search TUI or plain CLI commands.

## Features

- **Daemon** — polls the clipboard every 500 ms and persists entries to `~/.goclip/history.json`
- **Interactive TUI** — scrollable picker with live fuzzy search, built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Non-interactive CLI** — `list`, `copy <id>`, and `clear` for scripting
- **Keeps up to 100 entries**, de-duplicating consecutive identical clips

## Requirements

- Go 1.21+
- macOS, Linux, or Windows (clipboard access via [`atotto/clipboard`](https://github.com/atotto/clipboard))

## Installation

```bash
git clone https://github.com/yourusername/goclip.git
cd goclip
go build -o goclip .
sudo mv goclip /usr/local/bin/   # optional: put it on PATH
```

Or install directly with `go install`:

```bash
go install github.com/yourusername/goclip@latest
```

## Usage

### Typical workflow

```bash
# 1. Start the daemon (keep this running in a background tab or as a service)
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
| `goclip clear` | Wipe all saved history |
| `goclip help` | Show usage |

### TUI keybindings

| Key | Action |
|---|---|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Copy selected entry |
| `/` | Enter search mode |
| `Esc` | Exit search / quit |
| `q` | Quit |

## Storage

History is stored in `~/.goclip/history.json`. The file is plain JSON — you can inspect, back up, or edit it manually.

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

### Linux (systemd user service)

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
│   └── list.go        # list, copy, clear subcommands
├── storage/
│   └── storage.go     # Read/write history.json
└── ui/
    └── tui.go         # Bubbletea TUI picker
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT
