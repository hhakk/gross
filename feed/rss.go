package feed

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type RSSItem struct {
	url          string
	XTitle       string      `xml:"title"`
	XDescription HTMLContent `xml:"description"`
	XLink        string      `xml:"link"`
	XRead        bool        `xml:"read"`
}

type RSSChannel struct {
	XTitle       string     `xml:"title"`
	XDescription string     `xml:"description"`
	XLink        string     `xml:"link"`
	XItems       []*RSSItem `xml:"item"`
}

type RSS struct {
	url      string
	XMLName  xml.Name   `xml:"rss"`
	XChannel RSSChannel `xml:"channel"`
	AltTitle string
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
