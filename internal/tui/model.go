package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sivchari/crx/internal/registry"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170"))

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	tagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// Model represents the TUI state.
type Model struct {
	packages  []*registry.Package
	filtered  []*registry.Package
	selected  map[string]bool
	cursor    int
	textInput textinput.Model
	searching bool
	quitting  bool
	width     int
	height    int
}

// NewModel creates a new TUI model.
func NewModel(packages []*registry.Package) Model {
	ti := textinput.New()
	ti.Placeholder = "Search extensions..."
	ti.CharLimit = 50
	ti.Width = 40

	return Model{
		packages:  packages,
		filtered:  packages,
		selected:  make(map[string]bool),
		textInput: ti,
		searching: false,
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "/":
			m.searching = true
			m.textInput.Focus()
			return m, textinput.Blink

		case "esc":
			if m.searching {
				m.searching = false
				m.textInput.Blur()
				m.textInput.SetValue("")
				m.filtered = m.packages
			}
			return m, nil

		case "enter":
			if m.searching {
				m.searching = false
				m.textInput.Blur()
				return m, nil
			}
			// Toggle selection
			if len(m.filtered) > 0 {
				pkg := m.filtered[m.cursor]
				m.selected[pkg.Name] = !m.selected[pkg.Name]
			}
			return m, nil

		case "up", "k":
			if !m.searching && m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "j":
			if !m.searching && m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
			return m, nil

		case " ":
			if !m.searching && len(m.filtered) > 0 {
				pkg := m.filtered[m.cursor]
				m.selected[pkg.Name] = !m.selected[pkg.Name]
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	// Handle text input
	if m.searching {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		m.filterPackages()
		return m, cmd
	}

	return m, nil
}

func (m *Model) filterPackages() {
	query := strings.ToLower(m.textInput.Value())
	if query == "" {
		m.filtered = m.packages
		m.cursor = 0
		return
	}

	m.filtered = make([]*registry.Package, 0)
	for _, pkg := range m.packages {
		if strings.Contains(strings.ToLower(pkg.Name), query) ||
			strings.Contains(strings.ToLower(pkg.DisplayName), query) {
			m.filtered = append(m.filtered, pkg)
			continue
		}
		for _, tag := range pkg.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				m.filtered = append(m.filtered, pkg)
				break
			}
		}
	}
	m.cursor = 0
}

// View renders the UI.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("crx - Chrome Extension Manager"))
	b.WriteString("\n\n")

	// Search input
	if m.searching {
		b.WriteString(m.textInput.View())
	} else {
		b.WriteString(helpStyle.Render("Press / to search"))
	}
	b.WriteString("\n\n")

	// Package list
	for i, pkg := range m.filtered {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}

		checkbox := "[ ]"
		if m.selected[pkg.Name] {
			checkbox = "[x]"
		}

		line := fmt.Sprintf("%s%s %s", cursor, checkbox, pkg.DisplayName)

		if i == m.cursor {
			b.WriteString(selectedStyle.Render(line))
		} else {
			b.WriteString(normalStyle.Render(line))
		}
		b.WriteString("\n")

		// Show description for selected item
		if i == m.cursor && pkg.Description != "" {
			desc := fmt.Sprintf("      %s", pkg.Description)
			b.WriteString(tagStyle.Render(desc))
			b.WriteString("\n")
		}
	}

	// Help
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/k up • ↓/j down • space/enter select • / search • q quit"))
	b.WriteString("\n")

	// Selected count
	count := 0
	for _, v := range m.selected {
		if v {
			count++
		}
	}
	if count > 0 {
		b.WriteString(helpStyle.Render(fmt.Sprintf("\nSelected: %d extension(s)", count)))
	}

	return b.String()
}

// SelectedPackages returns the names of selected packages.
func (m Model) SelectedPackages() []string {
	var names []string
	for name, selected := range m.selected {
		if selected {
			names = append(names, name)
		}
	}
	return names
}

// Run starts the TUI.
func Run(packages []*registry.Package) ([]string, error) {
	m := NewModel(packages)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	return finalModel.(Model).SelectedPackages(), nil
}
