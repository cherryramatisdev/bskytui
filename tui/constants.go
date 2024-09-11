package tui

import "github.com/charmbracelet/lipgloss"

var (
	SpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#840806")).Bold(true)
)
