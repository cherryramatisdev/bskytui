package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
	"github.com/cherryramatisdev/bskytui/sdk"
)

type Post struct {
	ID      string
	Content string
	Summary string
	Author  string
	Langs   []string
}

func (p Post) Title() string {
	return p.Author
}

func (p Post) Description() string {
	return p.Content
}

func (p Post) FilterValue() string {
	return p.Author
}

type model struct {
	posts list.Model
}

func InitialModel(timeline sdk.Timeline) model {
	items := make([]list.Item, len(timeline.Feed))
	width, height, _ := term.GetSize(0)

	for i, post := range timeline.Feed {
		items[i] = Post{
			ID:      "id",
			Content: post.Post.Record.Text,
			Summary: "",
			Author:  post.Post.Author.DisplayName,
			Langs:   post.Post.Record.Langs,
		}
	}

	posts := list.New(items, list.NewDefaultDelegate(), width, height)
	posts.Title = "Timeline"

	return model{
		posts: posts,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.posts.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "ctrl-c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.posts, cmd = m.posts.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.posts.View()
}
