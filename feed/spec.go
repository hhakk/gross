package feed

import (
	"encoding/xml"

	"github.com/microcosm-cc/bluemonday"
)

var p = bluemonday.StrictPolicy()

type RSSItem struct {
	XTitle       string `xml:"title"`
	XDescription string `xml:"description"`
	XLink        string `xml:"link"`
}

type AtomEntry struct {
	XTitle   string `xml:"title"`
	XLink    string `xml:"link"`
	XSummary string `xml:"summary"`
	XUpdated string `xml:"updated"`
	XID      string `xml:"id"`
}

type RSSChannel struct {
	XTitle       string    `xml:"title"`
	XDescription string    `xml:"description"`
	XLink        string    `xml:"link"`
	XItems       []RSSItem `xml:"item"`
}

type Atom struct {
	XMLName  xml.Name    `xml:"feed"`
	XTitle   string      `xml:"title"`
	XLink    string      `xml:"link"`
	XUpdated string      `xml:"updated"`
	XID      string      `xml:"id"`
	XEntries []AtomEntry `xml:"entry"`
}

type RSS struct {
	XMLName  xml.Name   `xml:"rss"`
	XChannel RSSChannel `xml:"channel"`
}

type Item interface {
	FilterValue() string
	Title() string
	Description() string
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
	return i.XTitle
}

func (i RSSItem) Title() string {
	return i.XTitle
}

func (i RSSItem) Description() string {
	return i.XLink
}

func (i RSSItem) Link() string {
	return i.XLink
}

func (i RSSItem) Content() string {
	return p.Sanitize(i.XDescription)
}

func (a AtomEntry) FilterValue() string {
	return a.XTitle
}

func (a AtomEntry) Title() string {
	return a.XTitle
}

func (a AtomEntry) Link() string {
	return a.XLink
}

func (a AtomEntry) Description() string {
	return a.XLink
}

func (a AtomEntry) Content() string {
	return p.Sanitize(a.XSummary)
}

func (r RSS) FilterValue() string {
	return r.XChannel.XTitle
}

func (r RSS) Title() string {
	return r.XChannel.XTitle
}

func (r RSS) Description() string {
	return r.XChannel.XDescription
}

func (r RSS) Link() string {
	return r.XChannel.XLink
}

func (r RSS) Items() []Item {
	items := make([]Item, len(r.XChannel.XItems))
	for i, e := range r.XChannel.XItems {
		items[i] = e
	}
	return items
}

func (a Atom) FilterValue() string {
	return a.XTitle
}

func (a Atom) Title() string {
	return a.XTitle
}

func (a Atom) Description() string {
	return a.XTitle
}

func (a Atom) Link() string {
	return a.XLink
}

func (a Atom) Items() []Item {
	items := make([]Item, len(a.XEntries))
	for i, e := range a.XEntries {
		items[i] = e
	}
	return items
}
