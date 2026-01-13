package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/iml1s/antigravity-cleaner/internal/scanner"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	safeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	cautionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// DisplayScanResults shows scan results in a formatted table
func DisplayScanResults(items []scanner.CleanableItem) {
	if len(items) == 0 {
		fmt.Println("✨ No cleanable items found!")
		return
	}

	// Sort by size descending
	sort.Slice(items, func(i, j int) bool {
		return items[i].Size > items[j].Size
	})

	// Group by category
	categories := make(map[string][]scanner.CleanableItem)
	var totalSize int64

	for _, item := range items {
		categories[item.Category] = append(categories[item.Category], item)
		totalSize += item.Size
	}

	fmt.Println(titleStyle.Render("🔍 Scan Results"))
	fmt.Println()

	for category, catItems := range categories {
		var catSize int64
		for _, item := range catItems {
			catSize += item.Size
		}

		fmt.Printf("📁 %s (%s)\n", category, humanize.Bytes(uint64(catSize)))

		for _, item := range catItems {
			var levelStyle lipgloss.Style
			var levelIcon string
			switch item.SafeLevel {
			case "safe":
				levelStyle = safeStyle
				levelIcon = "✓"
			case "caution":
				levelStyle = cautionStyle
				levelIcon = "⚠"
			case "warning":
				levelStyle = warningStyle
				levelIcon = "⛔"
			}

			fmt.Printf("   %s %s %s\n",
				levelStyle.Render(levelIcon),
				levelStyle.Render(fmt.Sprintf("%-40s", item.Description)),
				humanize.Bytes(uint64(item.Size)))
		}
		fmt.Println()
	}

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📊 Total cleanable: %s\n", humanize.Bytes(uint64(totalSize)))
	fmt.Println()
	fmt.Println(helpStyle.Render("Legend: ✓ Safe  ⚠ Caution  ⛔ Warning"))
	fmt.Println(helpStyle.Render("Run 'agc clean' to interactively select items to clean"))
}

// SelectItems allows interactive selection of items to clean
func SelectItems(items []scanner.CleanableItem) []scanner.CleanableItem {
	if len(items) == 0 {
		return nil
	}

	// Sort by size descending
	sort.Slice(items, func(i, j int) bool {
		return items[i].Size > items[j].Size
	})

	m := initialModel(items)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return nil
	}

	fm := finalModel.(model)
	if fm.quit {
		return nil
	}

	var selected []scanner.CleanableItem
	for i, item := range items {
		if fm.selected[i] {
			selected = append(selected, item)
		}
	}

	return selected
}

// DisplayDryRun shows what would be cleaned without actually cleaning
func DisplayDryRun(items []scanner.CleanableItem) {
	var totalSize int64
	fmt.Println(titleStyle.Render("🔍 Dry Run - Would clean:"))
	fmt.Println()

	for _, item := range items {
		fmt.Printf("  • %s (%s)\n", item.Description, humanize.Bytes(uint64(item.Size)))
		fmt.Printf("    %s\n", item.Path)
		totalSize += item.Size
	}

	fmt.Println()
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📊 Would free: %s\n", humanize.Bytes(uint64(totalSize)))
}

// Interactive selection model using bubbletea
type model struct {
	items    []scanner.CleanableItem
	cursor   int
	selected map[int]bool
	quit     bool
	confirm  bool
}

func initialModel(items []scanner.CleanableItem) model {
	return model{
		items:    items,
		selected: make(map[int]bool),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case " ":
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "a":
			// Toggle all
			allSelected := true
			for i := range m.items {
				if !m.selected[i] {
					allSelected = false
					break
				}
			}
			for i := range m.items {
				m.selected[i] = !allSelected
			}
		case "s":
			// Select all safe items
			for i, item := range m.items {
				m.selected[i] = item.SafeLevel == "safe"
			}
		case "enter":
			hasSelection := false
			for _, v := range m.selected {
				if v {
					hasSelection = true
					break
				}
			}
			if hasSelection {
				m.confirm = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("🧹 Select items to clean"))
	s.WriteString("\n\n")

	var totalSelected int64
	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := "[ ]"
		if m.selected[i] {
			checked = "[x]"
			totalSelected += item.Size
		}

		var levelStyle lipgloss.Style
		switch item.SafeLevel {
		case "safe":
			levelStyle = safeStyle
		case "caution":
			levelStyle = cautionStyle
		case "warning":
			levelStyle = warningStyle
		}

		line := fmt.Sprintf("%s %s %s %s",
			cursor,
			checked,
			levelStyle.Render(fmt.Sprintf("%-45s", item.Description)),
			humanize.Bytes(uint64(item.Size)))

		if m.cursor == i {
			line = selectedStyle.Render(line)
		}

		s.WriteString(line + "\n")
	}

	s.WriteString("\n")
	s.WriteString(fmt.Sprintf("Selected: %s\n", humanize.Bytes(uint64(totalSelected))))
	s.WriteString("\n")
	s.WriteString(helpStyle.Render("↑/↓: Navigate • Space: Toggle • a: Toggle All • s: Select Safe • Enter: Confirm • q: Quit"))

	return s.String()
}
