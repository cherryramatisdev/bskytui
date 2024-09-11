package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cherryramatisdev/bskytui/tui"
)

func main() {
	p := tea.NewProgram(tui.NewApp())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
