package picker

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var ErrConfirmationCancelled = errors.New("confirmation cancelled")

var confirmStyles = struct {
	question lipgloss.Style
	warning  lipgloss.Style
	prompt   lipgloss.Style
	yes      lipgloss.Style
	no       lipgloss.Style
	hint     lipgloss.Style
	error_   lipgloss.Style
}{
	question: lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")).
		Bold(true),
	warning: lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")),
	prompt: lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true),
	yes: lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true),
	no: lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true),
	hint: lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")),
	error_: lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")),
}

type confirmModel struct {
	question  string
	confirmed bool
	cancelled bool
	focus     int
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		case "enter":
			if m.focus == 0 {
				m.confirmed = true
			}
			return m, tea.Quit
		case "left", "right", "tab", "shift+tab":
			m.focus = 1 - m.focus
			return m, nil
		case "y", "Y":
			m.confirmed = true
			return m, tea.Quit
		case "n", "N":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	var b strings.Builder

	b.WriteString(m.question)
	b.WriteString("\n\n")

	var yesLabel, noLabel string
	if m.focus == 0 {
		yesLabel = confirmStyles.yes.Render("▶ No ")
		noLabel = "  Yes "
	} else {
		yesLabel = "  No  "
		noLabel = confirmStyles.yes.Render("▶ Yes")
	}

	b.WriteString(yesLabel)
	b.WriteString("  ")
	b.WriteString(noLabel)
	b.WriteString("\n\n")

	b.WriteString(confirmStyles.hint.Render("←/→ or Tab to switch • Enter to confirm • n/N or Esc to cancel"))

	return b.String()
}

func (m confirmModel) Confirmed() bool {
	return m.confirmed
}

func (m confirmModel) Cancelled() bool {
	return m.cancelled
}

func Confirm(question string) (bool, error) {
	p := tea.NewProgram(confirmModel{question: question})
	model, err := p.Run()
	if err != nil {
		return false, fmt.Errorf("confirm error: %w", err)
	}

	m := model.(confirmModel)
	if m.Cancelled() {
		return false, ErrConfirmationCancelled
	}

	return m.Confirmed(), nil
}
