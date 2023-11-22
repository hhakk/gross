package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hhakk/gross/feed"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2, 1, 2)
var pageStyle = lipgloss.NewStyle().Width(20).Inherit(docStyle)

type sessionState uint

const (
	allFeedsView sessionState = iota
	singleFeedView
	singleItemView
)

type mainModel struct {
	URLs         []string
	fc           chan feed.FeedMessage
	state        sessionState
	allFeeds     list.Model
	singleFeed   list.Model
	singleItem   viewport.Model
	selectedFeed feed.Feed
	selectedItem feed.Item
}

type ListItem struct {
	title       string
	description string
}

func (l ListItem) FilterValue() string { return l.title }
func (l ListItem) Title() string       { return l.title }
func (l ListItem) Description() string { return l.description }

func formatContent(item feed.Item) string {
	return pageStyle.Render(fmt.Sprintf(
		`---%s---
Link: %s
---
%s`, item.Title(), item.Link(), item.Content()))
}

func newMainModel(urls []string) mainModel {
	w := 0
	h := 0
	fitems := make([]list.Item, len(urls))
	for i, url := range urls {
		fitems[i] = ListItem{title: url, description: "Loading..."}
	}
	allFeedList := list.New(fitems, list.NewDefaultDelegate(), w, h)
	allFeedList.Title = "gross | feeds"
	singleFeedList := list.New([]list.Item{}, list.NewDefaultDelegate(), w, h)
	return mainModel{
		URLs:       urls,
		fc:         make(chan feed.FeedMessage),
		state:      allFeedsView,
		allFeeds:   allFeedList,
		singleFeed: singleFeedList,
		singleItem: viewport.New(w, h),
	}
}

type feedLoadedMsg feed.FeedMessage

func receiveFeeds(c chan feed.FeedMessage) tea.Cmd {
	return func() tea.Msg {
		return feedLoadedMsg(<-c)
	}
}

func (m mainModel) Init() tea.Cmd {
	sender := func() tea.Msg {
		feed.GetFeeds(m.URLs, m.fc)
		return struct{}{}
	}
	return tea.Batch(sender, receiveFeeds(m.fc))
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
		// these keys quit the program
		case "ctrl+c", "q":
			return m, tea.Quit
			// these keys go back
		case "h", "right":
			switch m.state {
			case allFeedsView:
				return m, tea.Quit
			case singleFeedView:
				m.state = allFeedsView
			case singleItemView:
				m.state = singleFeedView
			}
			// these keys select items
		case "l", "left", "enter":
			switch m.state {
			case allFeedsView:
				f, ok := m.allFeeds.SelectedItem().(feed.Feed)
				if ok {
					m.singleFeed.Title = fmt.Sprintf("gross | %s", f.Title())
					m.selectedFeed = f
					listItems := make([]list.Item, len(f.Items()))
					for ix, li := range f.Items() {
						listItems[ix] = li
					}
					cmd = m.singleFeed.SetItems(listItems)
					cmds = append(cmds, cmd)
					m.state = singleFeedView
				}
			case singleFeedView:
				i, ok := m.singleFeed.SelectedItem().(feed.Item)
				if ok {
					m.selectedItem = i
					m.singleItem.SetContent(formatContent(i))
					m.state = singleItemView
				}
			}

		}
		// handle other keypresses accordingly
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
		// we receive a new feed
	case feedLoadedMsg:
		// we keep on listening
		cmds = append(cmds, receiveFeeds(m.fc))
		// on error, we show errored output
		if msg.Error != nil {
			cmd = m.allFeeds.SetItem(msg.Index, ListItem{title: "Error", description: fmt.Sprintf("Error: %s", msg.Error)})
		} else {
			// otherwise, we show the item
			cmd = m.allFeeds.SetItem(msg.Index, *msg.Feed)
		}
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	switch m.state {
	case allFeedsView:
		return m.allFeeds.View()
	case singleFeedView:
		return m.singleFeed.View()
	case singleItemView:
		return m.singleItem.View()
	}
	return ""
}

func Run(urls []string) error {
	m := newMainModel(urls)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
