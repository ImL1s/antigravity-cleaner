package scanner

import (
	"os"
	"path/filepath"
	"runtime"
)

// CleanableItem represents a directory or file that can be cleaned
type CleanableItem struct {
	Path        string
	Size        int64
	Category    string
	Description string
	SafeLevel   string // "safe", "caution", "warning"
}

// getHomeDir returns the user's home directory
func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

// getDirSize calculates the total size of a directory
func getDirSize(path string) int64 {
	var size int64
	_ = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

// exists checks if a path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ScanAll scans all supported categories
func ScanAll() []CleanableItem {
	var results []CleanableItem
	results = append(results, ScanAntigravity()...)
	results = append(results, ScanFlutter("")...)
	results = append(results, ScanXcode()...)
	results = append(results, ScanAndroid()...)
	results = append(results, ScanVSCode()...)
	return results
}

// ScanAntigravity scans for Google Antigravity IDE cleanable items
func ScanAntigravity() []CleanableItem {
	var results []CleanableItem
	home := getHomeDir()

	// Platform-specific paths
	var paths []struct {
		path        string
		category    string
		description string
		safeLevel   string
	}

	switch runtime.GOOS {
	case "darwin":
		paths = []struct {
			path        string
			category    string
			description string
			safeLevel   string
		}{
			// ~/.gemini/antigravity - 主要資料
			{filepath.Join(home, ".gemini", "antigravity", "browser_recordings"), "Antigravity", "Session recordings (screenshots)", "safe"},
			{filepath.Join(home, ".gemini", "antigravity", "conversations"), "Antigravity", "Conversation history", "caution"},
			{filepath.Join(home, ".gemini", "antigravity", "brain"), "Antigravity", "AI memory cache", "caution"},
			{filepath.Join(home, ".gemini", "antigravity", "implicit"), "Antigravity", "Implicit data cache", "safe"},

			// ~/Library/Application Support/Antigravity
			{filepath.Join(home, "Library", "Application Support", "Antigravity", "CachedData"), "Antigravity", "JS/WASM cached data", "safe"},
			{filepath.Join(home, "Library", "Application Support", "Antigravity", "Code Cache"), "Antigravity", "Code cache", "safe"},
			{filepath.Join(home, "Library", "Application Support", "Antigravity", "User", "_extensions-disabled"), "Antigravity", "Disabled extensions backup", "safe"},
			{filepath.Join(home, "Library", "Application Support", "Antigravity", "DawnWebGPUCache"), "Antigravity", "WebGPU cache", "safe"},
			{filepath.Join(home, "Library", "Application Support", "Antigravity", "DawnGraphiteCache"), "Antigravity", "Graphite cache", "safe"},
			{filepath.Join(home, "Library", "Application Support", "Antigravity", "User", "workspaceStorage"), "Antigravity", "Workspace storage", "caution"},

			// ~/.antigravity
			{filepath.Join(home, ".antigravity", "extensions"), "Antigravity", "Old extension versions", "safe"},
		}
	case "windows":
		appData := os.Getenv("APPDATA")
		localAppData := os.Getenv("LOCALAPPDATA")
		paths = []struct {
			path        string
			category    string
			description string
			safeLevel   string
		}{
			{filepath.Join(home, ".gemini", "antigravity", "browser_recordings"), "Antigravity", "Session recordings", "safe"},
			{filepath.Join(home, ".gemini", "antigravity", "conversations"), "Antigravity", "Conversation history", "caution"},
			{filepath.Join(home, ".gemini", "antigravity", "brain"), "Antigravity", "AI memory cache", "caution"},
			{filepath.Join(appData, "Antigravity", "CachedData"), "Antigravity", "Cached data", "safe"},
			{filepath.Join(appData, "Antigravity", "Code Cache"), "Antigravity", "Code cache", "safe"},
			{filepath.Join(localAppData, "Antigravity", "CachedData"), "Antigravity", "Local cached data", "safe"},
		}
	case "linux":
		configDir := filepath.Join(home, ".config")
		paths = []struct {
			path        string
			category    string
			description string
			safeLevel   string
		}{
			{filepath.Join(home, ".gemini", "antigravity", "browser_recordings"), "Antigravity", "Session recordings", "safe"},
			{filepath.Join(home, ".gemini", "antigravity", "conversations"), "Antigravity", "Conversation history", "caution"},
			{filepath.Join(home, ".gemini", "antigravity", "brain"), "Antigravity", "AI memory cache", "caution"},
			{filepath.Join(configDir, "Antigravity", "CachedData"), "Antigravity", "Cached data", "safe"},
			{filepath.Join(configDir, "Antigravity", "Code Cache"), "Antigravity", "Code cache", "safe"},
		}
	}

	for _, p := range paths {
		if exists(p.path) {
			size := getDirSize(p.path)
			if size > 0 {
				results = append(results, CleanableItem{
					Path:        p.path,
					Size:        size,
					Category:    p.category,
					Description: p.description,
					SafeLevel:   p.safeLevel,
				})
			}
		}
	}

	return results
}

// ScanFlutter scans for Flutter project build directories
func ScanFlutter(basePath string) []CleanableItem {
	var results []CleanableItem
	home := getHomeDir()

	if basePath == "" {
		basePath = filepath.Join(home, "Documents")
	}

	// Find Flutter projects by looking for pubspec.yaml
	_ = filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden directories and node_modules
		if info.IsDir() && (info.Name()[0] == '.' || info.Name() == "node_modules") {
			return filepath.SkipDir
		}

		// Look for build directories in Flutter projects
		if info.IsDir() && info.Name() == "build" {
			// Check if parent has pubspec.yaml (Flutter project)
			parent := filepath.Dir(path)
			if exists(filepath.Join(parent, "pubspec.yaml")) {
				size := getDirSize(path)
				if size > 100*1024*1024 { // Only show if > 100MB
					results = append(results, CleanableItem{
						Path:        path,
						Size:        size,
						Category:    "Flutter",
						Description: "Build directory: " + filepath.Base(parent),
						SafeLevel:   "safe",
					})
				}
			}
		}

		// Look for .dart_tool directories
		if info.IsDir() && info.Name() == ".dart_tool" {
			size := getDirSize(path)
			if size > 50*1024*1024 { // Only show if > 50MB
				parent := filepath.Dir(path)
				results = append(results, CleanableItem{
					Path:        path,
					Size:        size,
					Category:    "Flutter",
					Description: ".dart_tool: " + filepath.Base(parent),
					SafeLevel:   "safe",
				})
			}
		}

		return nil
	})

	// Global Dart/Flutter caches
	pubCache := filepath.Join(home, ".pub-cache")
	if exists(pubCache) {
		size := getDirSize(pubCache)
		if size > 100*1024*1024 {
			results = append(results, CleanableItem{
				Path:        pubCache,
				Size:        size,
				Category:    "Flutter",
				Description: "Pub package cache",
				SafeLevel:   "caution",
			})
		}
	}

	return results
}

// ScanXcode scans for Xcode cleanable items (macOS only)
func ScanXcode() []CleanableItem {
	var results []CleanableItem

	if runtime.GOOS != "darwin" {
		return results
	}

	home := getHomeDir()

	paths := []struct {
		path        string
		description string
		safeLevel   string
	}{
		{filepath.Join(home, "Library", "Developer", "Xcode", "DerivedData"), "Xcode DerivedData", "safe"},
		{filepath.Join(home, "Library", "Developer", "Xcode", "iOS DeviceSupport"), "iOS DeviceSupport", "safe"},
		{filepath.Join(home, "Library", "Developer", "Xcode", "watchOS DeviceSupport"), "watchOS DeviceSupport", "safe"},
		{filepath.Join(home, "Library", "Developer", "Xcode", "Archives"), "Xcode Archives", "caution"},
		{filepath.Join(home, "Library", "Developer", "CoreSimulator", "Caches"), "Simulator Caches", "safe"},
	}

	for _, p := range paths {
		if exists(p.path) {
			size := getDirSize(p.path)
			if size > 100*1024*1024 { // Only show if > 100MB
				results = append(results, CleanableItem{
					Path:        p.path,
					Size:        size,
					Category:    "Xcode",
					Description: p.description,
					SafeLevel:   p.safeLevel,
				})
			}
		}
	}

	return results
}

// ScanAndroid scans for Android development cleanable items
func ScanAndroid() []CleanableItem {
	var results []CleanableItem
	home := getHomeDir()

	paths := []struct {
		path        string
		description string
		safeLevel   string
	}{
		{filepath.Join(home, ".gradle", "caches"), "Gradle caches", "safe"},
		{filepath.Join(home, ".gradle", "wrapper", "dists"), "Gradle distributions", "caution"},
		{filepath.Join(home, ".android", "cache"), "Android SDK cache", "safe"},
	}

	for _, p := range paths {
		if exists(p.path) {
			size := getDirSize(p.path)
			if size > 100*1024*1024 {
				results = append(results, CleanableItem{
					Path:        p.path,
					Size:        size,
					Category:    "Android",
					Description: p.description,
					SafeLevel:   p.safeLevel,
				})
			}
		}
	}

	// Scan for AVD images
	avdPath := filepath.Join(home, ".android", "avd")
	if exists(avdPath) {
		entries, err := os.ReadDir(avdPath)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() && filepath.Ext(entry.Name()) == ".avd" {
					fullPath := filepath.Join(avdPath, entry.Name())
					size := getDirSize(fullPath)
					if size > 1024*1024*1024 { // Only show if > 1GB
						results = append(results, CleanableItem{
							Path:        fullPath,
							Size:        size,
							Category:    "Android",
							Description: "AVD: " + entry.Name(),
							SafeLevel:   "warning",
						})
					}
				}
			}
		}
	}

	return results
}

// ScanVSCode scans for VS Code and variants cleanable items
func ScanVSCode() []CleanableItem {
	var results []CleanableItem
	home := getHomeDir()

	var basePaths []string

	switch runtime.GOOS {
	case "darwin":
		basePaths = []string{
			filepath.Join(home, "Library", "Application Support", "Code"),
			filepath.Join(home, "Library", "Application Support", "Cursor"),
		}
	case "windows":
		appData := os.Getenv("APPDATA")
		basePaths = []string{
			filepath.Join(appData, "Code"),
			filepath.Join(appData, "Cursor"),
		}
	case "linux":
		configDir := filepath.Join(home, ".config")
		basePaths = []string{
			filepath.Join(configDir, "Code"),
			filepath.Join(configDir, "Cursor"),
		}
	}

	cacheSubdirs := []string{"CachedData", "Code Cache", "CachedExtensions", "CachedExtensionVSIXs"}

	for _, basePath := range basePaths {
		for _, subdir := range cacheSubdirs {
			path := filepath.Join(basePath, subdir)
			if exists(path) {
				size := getDirSize(path)
				if size > 50*1024*1024 {
					results = append(results, CleanableItem{
						Path:        path,
						Size:        size,
						Category:    "VS Code",
						Description: filepath.Base(basePath) + " " + subdir,
						SafeLevel:   "safe",
					})
				}
			}
		}
	}

	return results
}

// ScanSimulators scans for old simulator runtimes
func ScanSimulators() []CleanableItem {
	// This is handled specially via xcrun simctl
	return []CleanableItem{}
}
