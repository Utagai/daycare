package cmd

import (
	"charm.land/lipgloss/v2"
)

var (
	HeaderColor = lipgloss.Color("121")
	BorderColor = lipgloss.Color("77")
	EvenRowFg   = lipgloss.Color("245")
	OddRowFg    = lipgloss.Color("250")

	SuccessColor = lipgloss.Color("10")
	WarningColor = lipgloss.Color("11")
	ErrorColor   = lipgloss.Color("9")
	MutedColor   = lipgloss.Color("243")
)

var (
	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	PromptStyle = lipgloss.NewStyle().
			Foreground(HeaderColor).
			Bold(true)

	MutedStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	ProjectNameStyle = lipgloss.NewStyle().
				Foreground(HeaderColor).
				Bold(true)

	PathStyle = lipgloss.NewStyle().
			Foreground(MutedColor)
)
