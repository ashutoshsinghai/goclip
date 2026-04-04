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

# 2. Open the interactive picker whenever you need an old clip
goclip pick
```

### All commands

| Command | Description |
|---|---|
| `goclip daemon` | Start clipboard watcher in background |
| `goclip stop` | Stop background daemon |
| `goclip status` | Show whether daemon is running |
| `goclip run` | Run clipboard watcher in foreground |
| `goclip pick` | Open the interactive TUI picker |
| `goclip list` | Print history as a plain text table |
| `goclip search <keyword>` | Search history by keyword |
| `goclip copy <id>` | Put a historical entry back on the clipboard |
| `goclip pin <id>` | Pin/unpin an entry (pinned items stay at the top) |
| `goclip clear` | Wipe all saved history |
| `goclip upgrade` | Upgrade to the latest version |
| `goclip version` | Show current version |
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
│   ├── list.go        # list, search, copy, pin, clear subcommands
│   ├── upgrade.go     # Self-upgrade from GitHub Releases
│   └── extract.go     # tar.gz / zip extraction helpers
├── storage/
│   └── storage.go     # Read/write history.json
├── ui/
│   └── tui.go         # Bubbletea TUI picker
└── scripts/
    ├── install.sh     # Installer for macOS/Linux
    └── install.ps1    # Installer for Windows (PowerShell)
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — [Ashutosh Singhai](https://github.com/ashutoshsinghai)
