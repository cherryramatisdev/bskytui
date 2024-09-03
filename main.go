package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cherryramatisdev/bskytui/sdk"
	"github.com/cherryramatisdev/bskytui/tui"
	"github.com/cherryramatisdev/bskytui/util"
	"golang.org/x/term"
)

func main() {

	fmt.Println("Username (example.bsky.social):")
	fmt.Print("> @")

	var user string
	fmt.Scanln(&user)
	if strings.HasPrefix(user, "@") {
		user = user[1:]
	}

	fmt.Println("Password:")
	fmt.Print("> ")
	bPass, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	pass := string(bPass)

	ctx, err := sdk.Authenticate(sdk.AuthUser{
		Identifier: user,
		Password:   pass,
	})

	// TODO: remove this panic for a prettier handle (when it's TUI)
	if err != nil {
		panic(err)
	}

	if util.IsDebug() {
		_, _ = sdk.GetTimeline(ctx)
		return
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
