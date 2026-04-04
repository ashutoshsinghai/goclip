# Contributing to goclip

Thanks for your interest in contributing. Here is everything you need to get started.

## Prerequisites

- Go 1.21+
- A working clipboard (macOS, Linux with `xclip`/`xdotool`, or Windows)

## Getting the code

```bash
git clone https://github.com/yourusername/goclip.git
cd goclip
go mod download
```

## Building and running locally

```bash
go build -o goclip .
./goclip help
```

Run the daemon in one terminal and the picker in another to exercise the full flow:

```bash
./goclip daemon      # terminal 1
./goclip pick        # terminal 2
```

## Running tests

```bash
go test ./...
```

There are currently no test files — adding them is a great first contribution.

## Code layout

| Path | Responsibility |
|---|---|
| `main.go` | Argument parsing and routing only |
| `commands/daemon.go` | Clipboard polling loop |
| `commands/list.go` | `list`, `copy`, `clear` subcommands |
| `storage/storage.go` | Load/save `~/.goclip/history.json` |
| `ui/tui.go` | Bubbletea TUI model and view |

Keep new logic in the appropriate package. `main.go` should stay thin.

## Making changes

1. Fork the repository and create a feature branch:

   ```bash
   git checkout -b feat/your-feature
   ```

2. Make your changes. Keep commits small and focused.

3. Ensure the project still builds and passes `go vet`:

   ```bash
   go build ./...
   go vet ./...
   ```

4. Open a pull request against `main`. Describe what you changed and why.

## Pull request guidelines

- **One concern per PR.** Bug fixes, features, and refactors should be separate PRs.
- **Describe the motivation** in the PR body — what problem does this solve?
- **Keep diffs readable.** Avoid mixing formatting changes with logic changes.
- PRs that add or improve test coverage are especially welcome.

## Reporting bugs

Open a GitHub issue with:

1. Your OS and Go version (`go version`)
2. Steps to reproduce
3. What you expected vs. what happened

## Suggesting features

Open a GitHub issue with a clear description of the use case. Starting a discussion before writing code avoids wasted effort.

## Code style

Follow standard Go conventions (`gofmt`, `go vet`). There is no linter config yet — just match the style of the surrounding code.

## License

By contributing you agree that your changes will be licensed under the project's [MIT License](LICENSE).
