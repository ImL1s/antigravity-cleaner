package main

import (
	"fmt"
	"os"

	"github.com/iml1s/antigravity-cleaner/internal/cleaner"
	"github.com/iml1s/antigravity-cleaner/internal/scanner"
	"github.com/iml1s/antigravity-cleaner/internal/ui"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "agc",
		Short: "Antigravity Cleaner - Clean up IDE caches and build artifacts",
		Long: `Antigravity Cleaner (agc) is a cross-platform CLI tool for cleaning up
development environment caches, build artifacts, and other temporary files.

Supports:
  - Google Antigravity IDE (session recordings, conversations, caches)
  - Flutter/Dart (build directories, .dart_tool)
  - Xcode (DerivedData, iOS DeviceSupport)
  - Android Studio (.gradle caches, AVD images)
  - VS Code and variants (CachedData, extensions)
  - iOS/Android Simulators (old runtimes)`,
		Version: version,
	}

	// Scan command
	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "Scan for cleanable items",
		Long:  "Scan your system for IDE caches, build artifacts, and other cleanable items.",
		Run: func(cmd *cobra.Command, args []string) {
			results := scanner.ScanAll()
			ui.DisplayScanResults(results)
		},
	}

	// Clean command
	var cleanAll bool
	var cleanDryRun bool
	var cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Clean up caches and build artifacts",
		Long:  "Interactively select and clean IDE caches, build artifacts, and temporary files.",
		Run: func(cmd *cobra.Command, args []string) {
			results := scanner.ScanAll()
			if len(results) == 0 {
				fmt.Println("No cleanable items found.")
				return
			}

			var toClean []scanner.CleanableItem
			if cleanAll {
				toClean = results
			} else {
				toClean = ui.SelectItems(results)
			}

			if len(toClean) == 0 {
				fmt.Println("No items selected for cleaning.")
				return
			}

			if cleanDryRun {
				ui.DisplayDryRun(toClean)
				return
			}

			cleaner.CleanItems(toClean)
		},
	}
	cleanCmd.Flags().BoolVarP(&cleanAll, "all", "a", false, "Clean all items without prompting")
	cleanCmd.Flags().BoolVarP(&cleanDryRun, "dry-run", "n", false, "Show what would be cleaned without actually cleaning")

	// Antigravity-specific command
	var agCmd = &cobra.Command{
		Use:   "antigravity",
		Short: "Clean Antigravity IDE specific caches",
		Long:  "Clean Google Antigravity IDE session recordings, conversations, and caches.",
		Run: func(cmd *cobra.Command, args []string) {
			results := scanner.ScanAntigravity()
			if len(results) == 0 {
				fmt.Println("No Antigravity cleanable items found.")
				return
			}
			toClean := ui.SelectItems(results)
			if len(toClean) > 0 {
				cleaner.CleanItems(toClean)
			}
		},
	}

	// Flutter-specific command
	var flutterPath string
	var flutterCmd = &cobra.Command{
		Use:   "flutter",
		Short: "Clean Flutter project build directories",
		Long:  "Scan and clean Flutter project build directories, .dart_tool, and Pods.",
		Run: func(cmd *cobra.Command, args []string) {
			results := scanner.ScanFlutter(flutterPath)
			if len(results) == 0 {
				fmt.Println("No Flutter cleanable items found.")
				return
			}
			toClean := ui.SelectItems(results)
			if len(toClean) > 0 {
				cleaner.CleanItems(toClean)
			}
		},
	}
	flutterCmd.Flags().StringVarP(&flutterPath, "path", "p", "", "Path to scan for Flutter projects (default: ~/Documents)")

	// Xcode-specific command
	var xcodeCmd = &cobra.Command{
		Use:   "xcode",
		Short: "Clean Xcode caches and derived data",
		Long:  "Clean Xcode DerivedData, iOS DeviceSupport, and archives.",
		Run: func(cmd *cobra.Command, args []string) {
			results := scanner.ScanXcode()
			if len(results) == 0 {
				fmt.Println("No Xcode cleanable items found.")
				return
			}
			toClean := ui.SelectItems(results)
			if len(toClean) > 0 {
				cleaner.CleanItems(toClean)
			}
		},
	}

	// Simulator command
	var simCmd = &cobra.Command{
		Use:   "simulator",
		Short: "Clean old simulator runtimes",
		Long:  "Remove unavailable iOS/watchOS/tvOS simulator runtimes and devices.",
		Run: func(cmd *cobra.Command, args []string) {
			results := scanner.ScanSimulators()
			if len(results) == 0 {
				fmt.Println("No old simulator runtimes found.")
				return
			}
			toClean := ui.SelectItems(results)
			if len(toClean) > 0 {
				cleaner.CleanSimulators(toClean)
			}
		},
	}

	rootCmd.AddCommand(scanCmd, cleanCmd, agCmd, flutterCmd, xcodeCmd, simCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
