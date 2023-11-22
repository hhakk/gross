package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hhakk/gross/feed"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type sessionState uint

const (
	allFeedsView sessionState = iota
	singleFeedView
	singleItemView
)

type mainModel struct {
	URLs       []string
	state      sessionState
	allFeeds   list.Model
	singleFeed list.Model
	singleItem viewport.Model
	feeds      []feed.Feed // interface
	items      []feed.Item // interface
    selectedFeed feed.Feed
    selectedItem feed.Item
	content    string
}

type feedLoadedMsg feed.Feed
type switchStateMsg sessionState

func loadFeeds(urls []string) tea.Msg {
	c := make(chan feed.Feed)
	feed.GetFeeds(urls, c)
	f := <-c
	return feedLoadedMsg(f)
}

func (m mainModel) Init() tea.Cmd {
	return func() { loadFeeds(m.URLs) }
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.allFeeds.SetSize(msg.Width-h, msg.Height-v)
		m.singleFeed.SetSize(msg.Width-h, msg.Height-v)
		m.singleItem.Width = msg.Width - h
		m.singleItem.Height = msg.Height - v
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
            return m, tea.Quit
        case "enter" {
            switch m.state {
            case allFeedsView:
                f, ok := m.allFeeds.SelectedItem().(feed.Feed)
                if ok {
                    m.selectedFeed = f
                }
            case singleFeedView:
                i, ok := m.singleFeed.SelectedItem().(feed.Item)
                if ok {
                    m.selectedItem = f
                }
            }

        }
        switch m.state {
        case allFeedsView:
            m.allFeeds, cmd = m.allFeeds.Update(msg)
            cmds = append(cmds, cmd)
        case singleFeedView:
            m.singleFeed, cmd = m.singleFeed.Update(msg)
            cmds = append(cmds, cmd)
        case singleItemView:
            m.singleItem, cmd = m.singleItem.Update(msg)
            cmds = append(cmds, cmd)
        }
	case feedLoadedMsg:
		m.feeds = append(m.feeds, msg)
	case switchStateMsg:
		m.state = msg
	}
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	switch m.state {
	case allFeedsView:
		m.allFeeds.View()
	case singleFeedView:
		m.singleFeed.View()
	case singleItemView:
		m.singleItem.View()
	}
}

func Run(urls []string) error {
	m := newMainModel(urls)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
