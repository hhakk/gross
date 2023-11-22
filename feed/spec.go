package feed

import "encoding/xml"

type RSSItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
}

type AtomEntry struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Summary string `xml:"summary"`
	Updated string `xml:"updated"`
	ID      string `xml:"id"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Link        string    `xml:"link"`
	Items       []RSSItem `xml:"item"`
}

type Atom struct {
	XMLName xml.Name    `xml:"feed"`
	Title   string      `xml:"title"`
	Link    string      `xml:"link"`
	Updated string      `xml:"updated"`
	ID      string      `xml:"id"`
	Entries []AtomEntry `xml:"entry"`
}

type RSS struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RSSChannel `xml:"channel"`
}

type Item interface {
	FilterValue() string
	Title() string
	Content() string
	Link() string
}

type Feed interface {
	FilterValue() string
	Title() string
	Description() string
	Link() string
	Items() []Item
}

func (i RSSItem) FilterValue() string {
	return i.Title
}

func (i RSSItem) Title() string {
	return i.Title
}

func (i RSSItem) Link() string {
	return i.Link
}

func (i RSSItem) Content() string {
	return i.Description
}

func (a AtomEntry) FilterValue() string {
	return a.Title
}

func (a AtomEntry) Title() string {
	return a.Title
}

func (a AtomEntry) Link() string {
	return a.Link
}

func (a AtomEntry) Content() string {
	return a.Summary
}

func (r RSS) FilterValue() string {
	return r.Channel.Title
}

func (r RSS) Title() string {
	return r.Channel.Title
}

func (r RSS) Description() string {
	return r.Channel.Description
}

func (r RSS) Link() string {
	return r.Channel.Link
}

func (r RSS) Items() string {
	return r.Channel.Items
}

func (a Atom) FilterValue() string {
	return a.Title
}

func (a Atom) Title() string {
	return a.Title
}

func (a Atom) Description() string {
	return a.Title
}

func (a Atom) Link() string {
	return a.Link
}

func (a Atom) Items() string {
	return a.Entries
}
