// Package picker provides an interactive fuzzy-search picker using Bubble Tea.
package picker

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ErrCancelled is returned when the user cancels the picker (presses ESC or Ctrl+C).
var ErrCancelled = errors.New("selection cancelled")

// Item represents an item that can be selected in the picker.
type Item interface {
	// FilterValue returns the string that will be used for filtering.
	FilterValue() string
	// DisplayValue returns the string that will be displayed in the list.
	DisplayValue() string
}

// Model represents the picker's state.
type Model struct {
	items        []Item
	filtered     []Item
	cursor       int
	textInput    textinput.Model
	selected     Item
	cancelled    bool
	width        int
	height       int
	title        string
	windowHeight int
	selectedIdx  int
}

// styles for the picker
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212"))
)

// New creates a new picker model with the given items and title.
// If selectedIdx >= 0 and < len(items), that item will be pre-selected.
func New(items []Item, title string, selectedIdx int) Model {
	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40

	cursor := 0
	if selectedIdx >= 0 && selectedIdx < len(items) {
		cursor = selectedIdx
	}

	return Model{
		items:        items,
		filtered:     items,
		cursor:       cursor,
		textInput:    ti,
		title:        title,
		windowHeight: 10,
		selectedIdx:  selectedIdx,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit

		case "enter":
			if len(m.filtered) > 0 {
				m.selected = m.filtered[m.cursor]
			}
			return m, tea.Quit

		case "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "ctrl+n":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
			return m, nil

		case "ctrl+u":
			m.textInput.SetValue("")
			m.filter()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Reserve lines for title, input, and status
		m.windowHeight = msg.Height - 5
		if m.windowHeight < 3 {
			m.windowHeight = 3
		}
	}

	// Handle text input
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	m.filter()

	return m, cmd
}

// filter updates the filtered list based on the current input.
func (m *Model) filter() {
	query := strings.ToLower(m.textInput.Value())
	if query == "" {
		m.filtered = m.items
	} else {
		var filtered []Item
		for _, item := range m.items {
			if fuzzyMatch(strings.ToLower(item.FilterValue()), query) {
				filtered = append(filtered, item)
			}
		}
		m.filtered = filtered
	}

	// Reset cursor if it's out of bounds
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// fuzzyMatch checks if the text contains all characters of the pattern in order.
func fuzzyMatch(text, pattern string) bool {
	patternIdx := 0
	for _, char := range text {
		if patternIdx < len(pattern) && char == rune(pattern[patternIdx]) {
			patternIdx++
		}
	}
	return patternIdx == len(pattern)
}

// View implements tea.Model.
func (m Model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render(m.title))
	b.WriteString("\n")

	// Text input
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	// Items list
	if len(m.filtered) == 0 {
		b.WriteString(dimStyle.Render("  No matches found"))
		b.WriteString("\n")
	} else {
		// Calculate visible window
		start := 0
		end := len(m.filtered)

		if len(m.filtered) > m.windowHeight {
			// Center the cursor in the visible window
			halfWindow := m.windowHeight / 2
			start = m.cursor - halfWindow
			if start < 0 {
				start = 0
			}
			end = start + m.windowHeight
			if end > len(m.filtered) {
				end = len(m.filtered)
				start = end - m.windowHeight
				if start < 0 {
					start = 0
				}
			}
		}

		for i := start; i < end; i++ {
			item := m.filtered[i]
			cursor := "  "
			style := normalStyle

			if i == m.cursor {
				cursor = cursorStyle.Render("> ")
				style = selectedStyle
			}

			b.WriteString(cursor)
			b.WriteString(style.Render(item.DisplayValue()))
			b.WriteString("\n")
		}
	}

	// Status line
	b.WriteString("\n")
	statusText := fmt.Sprintf("%d/%d", len(m.filtered), len(m.items))
	b.WriteString(dimStyle.Render(statusText))
	b.WriteString(dimStyle.Render("  • ↑/↓ navigate • enter select • esc cancel"))

	return b.String()
}

// Selected returns the selected item, or nil if cancelled.
func (m Model) Selected() Item {
	return m.selected
}

// Cancelled returns true if the user cancelled the selection.
func (m Model) Cancelled() bool {
	return m.cancelled
}

// Run runs the picker and returns the selected item.
// Returns ErrCancelled if the user cancels the selection.
// If selectedIdx >= 0, that item will be pre-selected in the list.
func Run(items []Item, title string, selectedIdx int) (Item, error) {
	if len(items) == 0 {
		return nil, errors.New("no items to select from")
	}

	p := tea.NewProgram(New(items, title, selectedIdx))
	model, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("picker error: %w", err)
	}

	m := model.(Model)
	if m.Cancelled() {
		return nil, ErrCancelled
	}

	if m.Selected() == nil {
		return nil, ErrCancelled
	}

	return m.Selected(), nil
}

// StringItem is a simple Item implementation for a string value.
type StringItem string

func (s StringItem) FilterValue() string  { return string(s) }
func (s StringItem) DisplayValue() string { return string(s) }

// ProjectItem wraps a project for use in the picker.
type ProjectItem struct {
	Name string
	Path string
}

func (p ProjectItem) FilterValue() string  { return p.Name }
func (p ProjectItem) DisplayValue() string { return p.Name }

// WorktreeItem wraps a worktree for use in the picker.
type WorktreeItem struct {
	Name   string
	Branch string
	Path   string
	IsMain bool
}

func (w WorktreeItem) FilterValue() string { return w.Name }
func (w WorktreeItem) DisplayValue() string {
	if w.IsMain {
		return w.Name + " (main)"
	}
	return w.Name
}
