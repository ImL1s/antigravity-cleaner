package cleaner

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/dustin/go-humanize"
	"github.com/iml1s/antigravity-cleaner/internal/scanner"
)

// CleanItems removes the specified cleanable items
func CleanItems(items []scanner.CleanableItem) {
	var totalCleaned int64
	var successCount, failCount int

	fmt.Println("\n🧹 Cleaning...")

	for _, item := range items {
		fmt.Printf("  Removing %s... ", item.Description)

		err := os.RemoveAll(item.Path)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			failCount++
		} else {
			fmt.Printf("✓ %s freed\n", humanize.Bytes(uint64(item.Size)))
			totalCleaned += item.Size
			successCount++
		}
	}

	fmt.Println()
	fmt.Printf("✨ Done! Cleaned %d items, freed %s\n", successCount, humanize.Bytes(uint64(totalCleaned)))
	if failCount > 0 {
		fmt.Printf("⚠️  %d items failed to clean\n", failCount)
	}
}

// CleanSimulators removes unavailable simulator runtimes and devices
func CleanSimulators(items []scanner.CleanableItem) {
	if runtime.GOOS != "darwin" {
		fmt.Println("Simulator cleanup is only available on macOS")
		return
	}

	fmt.Println("\n🧹 Cleaning simulators...")

	// Delete unavailable devices
	fmt.Print("  Removing unavailable devices... ")
	cmd := exec.Command("xcrun", "simctl", "delete", "unavailable")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✓")
	}

	// For runtime deletion, we need to handle each runtime specifically
	for _, item := range items {
		if item.Category == "Simulator" {
			fmt.Printf("  Removing %s... ", item.Description)
			cmd := exec.Command("xcrun", "simctl", "runtime", "delete", item.Path)
			err := cmd.Run()
			if err != nil {
				fmt.Printf("❌ Failed: %v\n", err)
			} else {
				fmt.Printf("✓ %s freed\n", humanize.Bytes(uint64(item.Size)))
			}
		}
	}

	fmt.Println("\n✨ Simulator cleanup complete!")
}

// CleanFlutterProject runs flutter clean in a project directory
func CleanFlutterProject(projectPath string) error {
	cmd := exec.Command("flutter", "clean")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
