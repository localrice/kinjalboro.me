package markdown

import (
	"html/template"
	"time"
)

type Post struct {
	Title   string
	Date    time.Time
	Slug    string
	Content template.HTML
}
