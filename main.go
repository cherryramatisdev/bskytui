package main

import (
	"os"

	"github.com/cherryramatisdev/bskytui/sdk"
	"github.com/cherryramatisdev/bskytui/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ctx, err := sdk.Authenticate(sdk.AuthUser{
		// TODO: Implement proper text inputs on the TUI to handle user and
		// password
		Identifier: os.Getenv("BSKY_USER"),
		Password:   os.Getenv("BSKY_PASSWORD"),
	})

	// TODO: remove this panic for a prettier handle (when it's TUI)
	if err != nil {
		panic(err)
	}

	timeline, err := sdk.GetTimeline(ctx)

	// TODO: remove this panic for a prettier handle (when it's TUI)
	if err != nil {
		panic(err)
	}

	p := tea.NewProgram(tui.InitialModel(timeline))

	// TODO: remove this panic for a prettier handle (when it's TUI)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
