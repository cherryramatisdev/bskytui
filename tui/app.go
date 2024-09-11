package tui

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cherryramatisdev/bskytui/sdk"
)

type App struct {
	Component
	err     error
	posts   []*Post
	spinner spinner.Model
}

func NewApp() *App {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return &App{
		err:     nil,
		posts:   []*Post{},
		spinner: s,
	}
}

func (app *App) Init() tea.Cmd {
	return tea.Batch(app.FetchTimeline, app.spinner.Tick)
}

func (app *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	component, componentCmd := app.Component.Update(msg)
	app.Component = *component.(*Component)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			app.quitting = true
			return app, tea.Quit
		}
	case *AuthAskToLogin:
		login := NewLogin(app)
		return login, login.Init()
	case *TimelineSuccess:
		app.loading = false
		return app, nil
	case *TimelineError:
		app.loading = false
		app.err = errors.New("Coudn't load the timeline correctly")
		app.quitting = true
		return app, tea.Quit
	}

	var spinnerCmd tea.Cmd
	app.spinner, spinnerCmd = app.spinner.Update(msg)

	return app, tea.Batch(componentCmd, spinnerCmd)
}

func (app *App) View() string {
	var content string

	if app.loading {
		content = SpinnerStyle.Render(fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", app.spinner.View()))
	}

	if app.err != nil {
		content += fmt.Sprintf("%s: %s", ErrorStyle.Render("ERROR"), app.err.Error())
	}

	if app.quitting {
		content += "\n"
	}

	return content
}

func (app *App) FetchTimeline() tea.Msg {
	var err error

	session, err := sdk.LoadAuthInfo()

	if err != nil {
		return &AuthAskToLogin{}
	}

	app.loading = true

	_, err = sdk.GetTimeline(&session)

	if err != nil {
		return &TimelineError{}
	}

	// feed := slices.DeleteFunc(timeline.Feed, func(el sdk.Feed) bool {
	// 	return el.Post.Record.Reply != nil
	// })

	return &TimelineSuccess{}
}
