package tui

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hhakk/gross/feed"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
	"github.com/spf13/viper"
)

var docStyle = lipgloss.NewStyle().Margin(4, 8)

type sessionState uint

const (
	allFeedsView sessionState = iota
	singleFeedView
	singleItemView
)

type mainModel struct {
	URLs         []feed.FeedSpec
	fc           chan feed.FeedMessage
	state        sessionState
	allFeeds     list.Model
	singleFeed   list.Model
	singleItem   viewport.Model
	selectedFeed feed.Feed
	selectedItem feed.Item
	contentReady bool
	width        int
	showRead     bool
}

type ListItem struct {
	title       string
	description string
}

func (l ListItem) FilterValue() string { return l.title }
func (l ListItem) Title() string       { return l.title }
func (l ListItem) Description() string { return l.description }

var unreadStyle = lipgloss.NewStyle().Inline(true).Bold(true).Foreground(lipgloss.Color("#fe4d93"))

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	styles := list.NewDefaultItemStyles()
	tR := styles.NormalTitle.Render
	dR := styles.NormalDesc.Render
	if index == m.Index() {
		tR = styles.SelectedTitle.Render
		dR = styles.SelectedDesc.Render
	}
	switch li := listItem.(type) {
	case ListItem:
		title := li.Title()
		desc := li.Description()
		fmt.Fprintf(w, "%s\n%s", tR(title), dR(desc))
	case feed.Feed:
		title := li.Title()
		desc := li.Description()
		n := len(li.Items())
		ur := 0
		for _, i := range li.Items() {
			if !i.IsRead() {
				ur += 1
			}
		}
		unread := unreadStyle.Render(fmt.Sprintf("(%d/%d)", ur, n))
		title = fmt.Sprintf("%s %s", unread, title)
		fmt.Fprintf(w, "%s\n%s", tR(title), dR(desc))
	case feed.Item:
		title := li.Title()
		desc := li.Description()
		if !li.IsRead() {
			title = unreadStyle.Render("(N) ") + title
		}
		fmt.Fprintf(w, "%s\n%s", tR(title), dR(desc))
	default:
		return
	}
}

func newMainModel(urls []feed.FeedSpec) mainModel {
	w := 0
	h := 0
	fitems := make([]list.Item, len(urls))
	for i, url := range urls {
		title := url.URL
		if url.AltName != "" {
			title = url.AltName
		}
		fitems[i] = ListItem{title: title, description: "Loading..."}
	}
	allFeedList := list.New(fitems, itemDelegate{}, w, h)
	allFeedList.Title = "gross | feeds"
	singleFeedList := list.New([]list.Item{}, itemDelegate{}, w, h)
	return mainModel{
		URLs:         urls,
		fc:           make(chan feed.FeedMessage),
		state:        allFeedsView,
		allFeeds:     allFeedList,
		singleFeed:   singleFeedList,
		singleItem:   viewport.New(w, h),
		selectedItem: nil,
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

func formatContent(item feed.Item, width int) string {
	ls := list.DefaultStyles()
	return fmt.Sprintf("%s\n\n%s\n\n%s",
		ls.Title.Render(fmt.Sprintf("gross | %s", item.Title())),
		ls.StatusBar.Render(
			wrap.String(
				fmt.Sprintf("Link: %s", item.Link()),
				width,
			),
		),
		wordwrap.String(item.Content(), width),
	)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	width := 0
	height := 0
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		width = msg.Width - h
		m.width = width
		height = msg.Height - v
		m.allFeeds.SetSize(width, height)
		m.singleFeed.SetSize(width, height)
		m.singleItem.Width = width
		m.singleItem.Height = height
		if m.selectedItem != nil {
			m.singleItem.SetContent(formatContent(m.selectedItem, m.width))
		}
	case tea.KeyMsg:
		switch msg.String() {
		// these keys quit the program
		case "ctrl+c", "q":
			for _, i := range m.allFeeds.Items() {
				f, ok := i.(feed.Feed)
				if ok {
					feed.SaveFeed(f)
				}
			}
			return m, tea.Quit
			// these keys go back
		case "h", "left":
			switch m.state {
			case allFeedsView:
				for _, i := range m.allFeeds.Items() {
					f, ok := i.(feed.Feed)
					if ok {
						feed.SaveFeed(f)
					}
				}
				return m, tea.Quit
			case singleFeedView:
				m.state = allFeedsView
			case singleItemView:
				m.state = singleFeedView
			}
		case "r":
			// toggle visible items
			if m.state == singleFeedView {
				m.showRead = !m.showRead
				f := m.selectedFeed
				listItems := make([]list.Item, 0)
				for _, li := range f.Items() {
					if m.showRead && li.IsRead() {
						continue
					}
					listItems = append(listItems, li)
				}
				cmd = m.singleFeed.SetItems(listItems)
				cmds = append(cmds, cmd)
			}

		case "A":
			if m.state == singleFeedView {
				for _, i := range m.selectedFeed.Items() {
					i.SetRead(true)
				}
			}
		case "a":
			// toggle read state
			i, ok := m.singleFeed.SelectedItem().(feed.Item)
			if m.state == singleFeedView && ok {
				i.SetRead(!i.IsRead())
			}
		case "l", "right", "enter":
			// these keys select items
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
					i.SetRead(true)
					m.selectedItem = i
					m.singleItem.SetContent(formatContent(i, m.width))
					m.state = singleItemView
				}
			case singleItemView:
				if m.selectedItem != nil {
					scmd := exec.Command(viper.GetString("browsercmd"), m.selectedItem.Link())
					err := scmd.Run()
					if err != nil {
						m.singleItem.SetContent(wrap.String(fmt.Sprintf("%s\n", err), m.width))
					}

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
			cmd = m.allFeeds.SetItem(
				msg.Index,
				ListItem{
					title:       msg.URL,
					description: fmt.Sprintf("Error: %s", msg.Error),
				},
			)
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

func Run(urls []feed.FeedSpec) error {
	m := newMainModel(urls)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
