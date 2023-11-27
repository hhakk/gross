package feed

import (
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
