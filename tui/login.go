package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/cherryramatisdev/bskytui/sdk"
)

type Login struct {
	Component
	parent   tea.Model
	form     *huh.Form
	username huh.Field
	password huh.Field
	err      string
	spinner  spinner.Model
}

func NewLogin(parent tea.Model) *Login {
	username := huh.NewInput().Key("username").Title("Username").Validate(required("Username"))
	password := huh.NewInput().Key("password").Title("Password").Validate(required("Password")).EchoMode(huh.EchoModePassword)

	s := spinner.New()
	s.Spinner = spinner.Dot

	return &Login{
		parent:   parent,
		form:     huh.NewForm(huh.NewGroup(username, password)),
		username: username,
		password: password,
		spinner:  s,
	}
}

func (l *Login) Init() tea.Cmd {
	l.form.SubmitCmd = l.DoLogin
	return l.form.Init()
}

func (l *Login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	component, cmd := l.Component.Update(msg)
	l.Component = *component.(*Component)

	switch msg := msg.(type) {
	case error:
		l.loading = false
		return l, l.HandleError(msg)
	case *AuthSuccess:
		return l.parent, l.parent.Init()
	}

	form, formCmd := l.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		l.form = f
	}

	return l, tea.Batch(cmd, formCmd)
}

func (l *Login) View() string {
	form := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1).
		Render(l.form.View())

	error := ErrorStyle.Render(l.err)

	if l.quitting {
		return ""
	}

	if l.loading {
		return lipgloss.Place(
			WindowWidth, WindowHeight,
			lipgloss.Center, 0.8,
			SpinnerStyle.Render(fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", l.spinner.View())),
		)
	}

	return lipgloss.Place(
		WindowWidth, WindowHeight,
		lipgloss.Center, 0.8,
		lipgloss.JoinVertical(lipgloss.Center, form, error),
	)
}

func (l *Login) DoLogin() tea.Msg {
	l.err = ""

	username := l.form.GetString("username")
	password := l.form.GetString("password")

	l.loading = true

	err := sdk.Authenticate(username, password)

	if err != nil {
		return err
	}

	return &AuthSuccess{}
}

func (l *Login) HandleError(err error) tea.Cmd {
	l.err = err.Error()
	l.form = huh.NewForm(huh.NewGroup(l.username, l.password))
	l.form.SubmitCmd = l.DoLogin
	cmd := l.form.Init()

	return cmd
}

func required(field string) func(value string) error {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required", field)
		}
		return nil
	}
}
