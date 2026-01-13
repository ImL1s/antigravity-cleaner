# Antigravity Cleaner

<p align="center">
  <strong>A cross-platform CLI tool for cleaning up development environment caches and build artifacts.</strong>
</p>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#supported-tools">Supported Tools</a> •
  <a href="#safety-levels">Safety Levels</a>
</p>

---

## Why?

As developers, our machines accumulate massive amounts of cache files, build artifacts, and temporary data from various IDEs and tools. This tool helps you:

- **Reclaim disk space** - Easily find and remove gigabytes of unnecessary files
- **Stay organized** - See exactly what's taking up space across all your dev tools
- **Stay safe** - Clear safety indicators help you avoid deleting important data

## Installation

### macOS (Homebrew)

```bash
brew install iml1s/tap/agc
```

### Windows (Scoop)

```powershell
scoop bucket add iml1s https://github.com/iml1s/scoop-bucket
scoop install agc
```

### Linux

**Debian/Ubuntu:**
```bash
# Download the latest .deb from GitHub releases
sudo dpkg -i agc_*_linux_amd64.deb
```

**RPM-based (Fedora, CentOS):**
```bash
sudo rpm -i agc_*_linux_amd64.rpm
```

### From Source

```bash
go install github.com/iml1s/antigravity-cleaner/cmd/agc@latest
```

### Manual Download

Download the latest binary from [GitHub Releases](https://github.com/iml1s/antigravity-cleaner/releases).

## Usage

### Scan for cleanable items

```bash
agc scan
```

Example output:
```
🔍 Scan Results

📁 Antigravity (17 GB)
   ✓ Session recordings (screenshots)         13 GB
   ✓ Disabled extensions backup               1.3 GB
   ⚠ AI memory cache                          810 MB
   ⚠ Conversation history                     808 MB

📁 Flutter (110 GB)
   ✓ Build directory: my_app                  8.5 GB
   ✓ Build directory: another_project         4.2 GB

📁 Xcode (5.6 GB)
   ✓ iOS DeviceSupport                        5.6 GB

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 Total cleanable: 155 GB
```

### Interactive cleanup

```bash
agc clean
```

Use arrow keys to navigate, space to select, and enter to confirm:
```
🧹 Select items to clean

> [x] ✓ Session recordings (screenshots)      13 GB
  [ ] ⚠ AI memory cache                       810 MB
  [x] ✓ Build directory: my_app               8.5 GB

Selected: 21.5 GB

↑/↓: Navigate • Space: Toggle • a: Toggle All • s: Select Safe • Enter: Confirm • q: Quit
```

### Other commands

```bash
# Clean all items without prompting
agc clean --all

# Preview what would be cleaned (dry run)
agc clean --dry-run

# Clean only Antigravity IDE caches
agc antigravity

# Clean only Flutter build directories
agc flutter
agc flutter --path ~/Projects  # Specify custom path

# Clean only Xcode caches
agc xcode

# Clean old simulator runtimes
agc simulator
```

## Supported Tools

### Google Antigravity IDE

**All Platforms:**
| Path | Description | Safety |
|------|-------------|:------:|
| `~/.gemini/antigravity/browser_recordings/` | Session replay screenshots | ✓ |
| `~/.gemini/antigravity/conversations/` | Conversation history | ⚠ |
| `~/.gemini/antigravity/brain/` | AI memory cache | ⚠ |
| `~/.gemini/antigravity/implicit/` | Implicit data cache | ✓ |

**macOS:**
| Path | Description | Safety |
|------|-------------|:------:|
| `~/Library/Application Support/Antigravity/CachedData/` | JS/WASM cache | ✓ |
| `~/Library/Application Support/Antigravity/Code Cache/` | Code cache | ✓ |
| `~/Library/Application Support/Antigravity/User/_extensions-disabled/` | Disabled extensions | ✓ |
| `~/Library/Application Support/Antigravity/DawnWebGPUCache/` | WebGPU cache | ✓ |
| `~/Library/Application Support/Antigravity/DawnGraphiteCache/` | Graphite cache | ✓ |
| `~/.antigravity/extensions/` | Old extension versions | ✓ |

**Windows:**
| Path | Description | Safety |
|------|-------------|:------:|
| `%APPDATA%\Antigravity\CachedData\` | Cached data | ✓ |
| `%APPDATA%\Antigravity\Code Cache\` | Code cache | ✓ |
| `%LOCALAPPDATA%\Antigravity\CachedData\` | Local cached data | ✓ |

**Linux:**
| Path | Description | Safety |
|------|-------------|:------:|
| `~/.config/Antigravity/CachedData/` | Cached data | ✓ |
| `~/.config/Antigravity/Code Cache/` | Code cache | ✓ |

### Flutter / Dart

**All Platforms:**
| Path | Description | Safety |
|------|-------------|:------:|
| `<project>/build/` | Build artifacts | ✓ |
| `<project>/.dart_tool/` | Dart tool cache | ✓ |
| `~/.pub-cache/` | Pub package cache | ⚠ |

### Xcode (macOS only)

| Path | Description | Safety |
|------|-------------|:------:|
| `~/Library/Developer/Xcode/DerivedData/` | Build cache | ✓ |
| `~/Library/Developer/Xcode/iOS DeviceSupport/` | Device symbols | ✓ |
| `~/Library/Developer/Xcode/watchOS DeviceSupport/` | Watch device symbols | ✓ |
| `~/Library/Developer/Xcode/Archives/` | App archives | ⚠ |
| `~/Library/Developer/CoreSimulator/Caches/` | Simulator caches | ✓ |

### Android Studio

**All Platforms:**
| Path | Description | Safety |
|------|-------------|:------:|
| `~/.gradle/caches/` | Gradle cache | ✓ |
| `~/.gradle/wrapper/dists/` | Gradle distributions | ⚠ |
| `~/.android/cache/` | Android SDK cache | ✓ |
| `~/.android/avd/*.avd/` | AVD images | ⛔ |

### VS Code & Variants (Cursor, etc.)

**macOS:**
| Path | Description | Safety |
|------|-------------|:------:|
| `~/Library/Application Support/Code/CachedData/` | Cached data | ✓ |
| `~/Library/Application Support/Code/Code Cache/` | Code cache | ✓ |
| `~/Library/Application Support/Cursor/CachedData/` | Cursor cached data | ✓ |

**Windows:**
| Path | Description | Safety |
|------|-------------|:------:|
| `%APPDATA%\Code\CachedData\` | Cached data | ✓ |
| `%APPDATA%\Code\Code Cache\` | Code cache | ✓ |
| `%APPDATA%\Cursor\CachedData\` | Cursor cached data | ✓ |

**Linux:**
| Path | Description | Safety |
|------|-------------|:------:|
| `~/.config/Code/CachedData/` | Cached data | ✓ |
| `~/.config/Code/Code Cache/` | Code cache | ✓ |
| `~/.config/Cursor/CachedData/` | Cursor cached data | ✓ |

## Safety Levels

| Icon | Level | Description |
|:----:|-------|-------------|
| ✓ | **Safe** | Can be deleted without any impact. Files will be regenerated automatically when needed. |
| ⚠ | **Caution** | May contain useful data. Review before deleting. Functionality won't break, but you may lose history or need to re-download. |
| ⛔ | **Warning** | Deleting may require significant reconfiguration or large re-downloads (e.g., AVD images). |

## Platform Support

| Platform | Status |
|----------|:------:|
| macOS (Intel) | ✅ |
| macOS (Apple Silicon) | ✅ |
| Windows (x64) | ✅ |
| Linux (x64) | ✅ |
| Linux (ARM64) | ✅ |

## Development

### Build from source

```bash
git clone https://github.com/iml1s/antigravity-cleaner.git
cd antigravity-cleaner
go build -o agc ./cmd/agc
```

### Run tests

```bash
go test ./...
```

### Create a release

```bash
# Install goreleaser
brew install goreleaser

# Create and push a tag
git tag v0.1.0
git push origin v0.1.0

# Release
goreleaser release --clean
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.

---

<p align="center">
  Made with ❤️ for developers who hate running out of disk space
</p>
