package feed

import (
	"encoding/xml"
	"fmt"
	"html"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var p = bluemonday.StrictPolicy()

type HTMLContent struct {
	Raw string `xml:",innerxml"`
}

func Escape(c HTMLContent) string {
	return strings.TrimSpace(
			p.Sanitize(
				html.UnescapeString(c.Raw),
			),
		)
}

type RSSItem struct {
	url          string
	XTitle       string      `xml:"title"`
	XDescription HTMLContent `xml:"description"`
	XLink        string      `xml:"link"`
	XRead        bool        `xml:"read"`
}

type AtomLink struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
}

type AtomEntry struct {
	url      string
	XTitle   string      `xml:"title"`
	XLink    AtomLink    `xml:"link`
	XSummary HTMLContent `xml:"summary"`
	XContent HTMLContent `xml:"content"`
	XUpdated string      `xml:"updated"`
	XD       string      `xml:"id"`
	XRead    bool        `xml:"read"`
}

type RSSChannel struct {
	XTitle       string     `xml:"title"`
	XDescription string     `xml:"description"`
	XLink        string     `xml:"link"`
	XItems       []*RSSItem `xml:"item"`
}

type Atom struct {
	url      string
	XMLName  xml.Name     `xml:"feed"`
	XTitle   string       `xml:"title"`
	XLink    AtomLink     `xml:"link"`
	XUpdated string       `xml:"updated"`
	XD       string       `xml:"id"`
	XEntries []*AtomEntry `xml:"entry"`
	AltTitle string
}

type RSS struct {
	url      string
	XMLName  xml.Name   `xml:"rss"`
	XChannel RSSChannel `xml:"channel"`
	AltTitle string
}

type Item interface {
	URL() string
	FilterValue() string
	Title() string
	Description() string
	Content() string
	Link() string
	IsRead() bool
	SetRead(read bool)
}

type Feed interface {
	URL() string
	FilterValue() string
	Title() string
	Description() string
	Link() string
	Items() []Item
	SetTitle(t string)
}

func (i *RSSItem) IsRead() bool {
	return i.XRead
}

func (i *RSSItem) SetRead(read bool) {
	i.XRead = read
}

func (i *RSSItem) URL() string {
	return i.url
}

func (i *RSSItem) FilterValue() string {
	return i.XTitle
}

func (i *RSSItem) Title() string {
	return i.XTitle
}

func (i *RSSItem) Description() string {
	return i.XLink
}

func (i *RSSItem) Link() string {
	link := i.XLink
	if len(link) > 1 && strings.HasPrefix(link, "/") {
		link = fmt.Sprintf("%s%s", i.url, link)
	}
	return link
}

func (i *RSSItem) Content() string {
	return Escape(i.XDescription)
}

func (a *AtomEntry) IsRead() bool {
	return a.XRead
}

func (a *AtomEntry) SetRead(read bool) {
	a.XRead = read
}

func (a *AtomEntry) URL() string {
	return a.url
}

func (a *AtomEntry) FilterValue() string {
	return a.XTitle
}

func (a *AtomEntry) Title() string {
	return a.XTitle
}

func (a *AtomEntry) Link() string {
	if a.XLink.Href != "" && a.XLink.Rel == "alternate" {
		link := a.XLink.Href
		if len(link) > 1 && strings.HasPrefix(link, "/") {
			link = fmt.Sprintf("%s%s", a.url, link)
		}
		return link
	}
	return a.XD
}

func (a *AtomEntry) Description() string {
	return a.Link()
}

func (a *AtomEntry) Content() string {
	if a.XContent.Raw != "" {
		return Escape(a.XContent)
	}
	return Escape(a.XSummary)
}

func (r *RSS) URL() string {
	return r.url
}

func (r *RSS) FilterValue() string {
	return r.XChannel.XTitle
}

func (r *RSS) SetTitle(t string) {
	r.AltTitle = t
}

func (r *RSS) Title() string {
	if r.AltTitle != "" {
		return r.AltTitle
	}
	return r.XChannel.XTitle
}

func (r *RSS) Description() string {
	return r.XChannel.XDescription
}

func (r *RSS) Link() string {
	return r.XChannel.XLink
}

func (r *RSS) Items() []Item {
	items := make([]Item, len(r.XChannel.XItems))
	for i, e := range r.XChannel.XItems {
		items[i] = e
	}
	return items
}

func (a *Atom) URL() string {
	return a.url
}

func (a *Atom) FilterValue() string {
	return a.XTitle
}

func (a *Atom) SetTitle(t string) {
	a.AltTitle = t
}

func (a *Atom) Title() string {
	if a.AltTitle != "" {
		return a.AltTitle
	}
	return a.XTitle
}

func (a *Atom) Description() string {
	return a.XTitle
}

func (a *Atom) Link() string {
	if a.XLink.Href != "" && a.XLink.Rel == "alternate" {
		return a.XLink.Href
	}
	return a.XD
}

func (a *Atom) Items() []Item {
	items := make([]Item, len(a.XEntries))
	for i, e := range a.XEntries {
		items[i] = e
	}
	return items
}
