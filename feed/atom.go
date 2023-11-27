package feed

import (
	"encoding/xml"
	"fmt"
	"strings"
)

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
