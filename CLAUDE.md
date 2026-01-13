# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
# Build
go build -o agc ./cmd/agc

# Run tests
go test ./...

# Run linter (used in CI)
golangci-lint run

# Build for all platforms (requires goreleaser)
goreleaser build --snapshot --clean

# Release (creates GitHub release + updates homebrew/scoop)
git tag v0.x.x && git push origin v0.x.x
# GitHub Actions will handle the rest
```

## Architecture

This is a Go CLI tool using cobra for commands and bubbletea for interactive UI.

```
cmd/agc/main.go      - CLI entry point, defines all commands (scan, clean, antigravity, flutter, xcode, simulator)
internal/
  scanner/           - Scans filesystem for cleanable items, platform-aware paths (darwin/windows/linux)
  cleaner/           - Executes cleanup operations (os.RemoveAll, xcrun simctl for simulators)
  ui/                - Terminal UI with lipgloss styling and bubbletea interactive selection
```

**Data Flow:** Commands call `scanner.Scan*()` → returns `[]CleanableItem` → `ui.SelectItems()` for interactive selection → `cleaner.CleanItems()` to delete.

**CleanableItem** has SafeLevel: "safe" (auto-regenerated), "caution" (may lose data), "warning" (requires significant re-download).

## Cross-Platform Paths

Scanner uses `runtime.GOOS` switch and `filepath.Join` for all paths. Key locations:
- **All platforms:** `~/.gemini/antigravity/`, `~/.gradle/`, `~/.pub-cache/`
- **macOS:** `~/Library/Application Support/`, `~/Library/Developer/`
- **Windows:** `%APPDATA%`, `%LOCALAPPDATA%`
- **Linux:** `~/.config/`

## Release Process

goreleaser handles cross-compilation and publishes to:
- GitHub Releases (binaries + deb/rpm)
- iml1s/homebrew-tap (Formula/agc.rb)
- iml1s/scoop-bucket (agc.json)

Secrets `HOMEBREW_TAP_GITHUB_TOKEN` and `SCOOP_GITHUB_TOKEN` must be set in repo settings.
