package tui

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
	"github.com/cherryramatisdev/bskytui/sdk"
)

type Author struct {
	DisplayName string
	Handle      string
}

type Post struct {
	ID      string
	Content string
	Summary string
	Author  Author
	Langs   []string
}

func (p Post) Title() string {
	return fmt.Sprintf("%s-%s", p.Author.DisplayName, p.Author.Handle)
}

func (p Post) Description() string {
	if p.Content != "" {
		return fmt.Sprintf("%s...", p.Content[:len(p.Content)/2])
	} else {
		return ""
	}
}

func (p Post) FilterValue() string {
	return p.Author.DisplayName
}

type TimelineModel struct {
	Component
	posts        list.Model
	selectedPost Post
}

func NewTimeline(timeline sdk.Timeline) *TimelineModel {
	// TODO: as soon as we're able to render replies properly, remove this
	feed := slices.DeleteFunc(timeline.Feed, func(el sdk.Feed) bool {
		return el.Post.Record.Reply != nil
	})

	items := make([]list.Item, len(feed))
	width, height, _ := term.GetSize(0)

	for i, post := range feed {
		items[i] = Post{
			ID:      "id",
			Content: post.Post.Record.Text,
			Summary: "",
			Author: Author{
				DisplayName: post.Post.Author.DisplayName,
				Handle:      post.Post.Author.Handle,
			},
			Langs: post.Post.Record.Langs,
		}
	}

	posts := list.New(items, list.NewDefaultDelegate(), width, height)
	posts.Title = "Timeline"

	return &TimelineModel{
		posts: posts,
	}
}

func (t *TimelineModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if t.posts.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "ctrl-c", "q":
			t.quitting = true
			return t, tea.Quit
		}
	}

	var cmd tea.Cmd
	t.posts, cmd = t.posts.Update(msg)
	return t, cmd
}

func (t *TimelineModel) View() string {
	if t.quitting {
		return ""
	}

	return t.posts.View()
}
