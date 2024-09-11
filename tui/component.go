package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

var (
	WindowHeight int
	WindowWidth  int
)

var _ tea.Model = new(Component)

type Component struct {
	quitting bool
	loading  bool
}

func (c *Component) Init() tea.Cmd {
	return nil
}

func (c *Component) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			c.quitting = true
			return c, tea.Quit
		}
	case tea.WindowSizeMsg:
		WindowWidth = msg.Width
		WindowHeight = msg.Height
		return c, tea.ClearScreen
	}
	return c, nil
}

func (c *Component) View() string {
	return ""
}
