package tui

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/x/term"
	"github.com/cherryramatisdev/bskytui/sdk"
	"github.com/cherryramatisdev/bskytui/util"
)

type Author struct {
	DisplayName string
	Handle      string
}

type Interactions struct {
	ReplyCount  int
	RepostCount int
	LikeCount   int
	QuoteCount  int
}

type Post struct {
	ID           string
	Content      string
	Summary      string
	Author       Author
	Langs        []string
	Interactions Interactions
}

func (p Post) Title() string {
	return fmt.Sprintf("%s (@%s) - (%s)", p.Author.DisplayName, p.Author.Handle, p.InteractionsString())
}

func (p Post) InteractionsString() string {
	return fmt.Sprintf("%v %s %v %s %v %s", p.Interactions.LikeCount, util.GetCountLabel(p.Interactions.LikeCount, "like"), p.Interactions.RepostCount+p.Interactions.QuoteCount, util.GetCountLabel(p.Interactions.RepostCount+p.Interactions.QuoteCount, "repost"), p.Interactions.ReplyCount, util.GetCountLabel(p.Interactions.ReplyCount, "reply"))
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

type view int

const (
	listView = iota
	detailView
)

type model struct {
	view         view
	posts        list.Model
	selectedPost Post
}

func InitialModel(timeline sdk.Timeline) model {
	feed := slices.DeleteFunc(timeline.Feed, func(el sdk.Feed) bool {
		return el.Post.Record.Reply != nil
	})

	items := make([]list.Item, len(feed))
	width, height, _ := term.GetSize(0)

	for i, c := range feed {
		items[i] = Post{
			ID:      "id",
			Content: c.Post.Record.Text,
			Summary: "",
			Author: Author{
				DisplayName: c.Post.Author.DisplayName,
				Handle:      c.Post.Author.Handle,
			},
			Interactions: Interactions{
				ReplyCount:  c.Post.ReplyCount,
				RepostCount: c.Post.RepostCount,
				LikeCount:   c.Post.LikeCount,
				QuoteCount:  c.Post.QuoteCount,
			},
			Langs: c.Post.Record.Langs,
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
		case "o", "enter":
			m.selectedPost = m.posts.SelectedItem().(Post)
			m.view = detailView
			return m, nil
		case "esc":
			if m.view == detailView {
				m.view = listView
			}

			if m.view == listView && m.posts.IsFiltered() {
				m.posts.ResetFilter()
			}

			return m, nil
		case "ctrl-c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd

	m.posts, cmd = m.posts.Update(msg)
	return m, cmd
}

func (m model) View() string {
	switch m.view {
	case detailView:
		view, err := glamour.Render(m.selectedPost.Content, "dark")
		// TODO: deal with errors in a better way
		if err != nil {
			panic(err)
		}
		return view
	default:
		return m.posts.View()
	}
}
