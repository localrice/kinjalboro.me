package markdown

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
)

type frontmatterData struct {
	Title string    `yaml:"title"`
	Date  time.Time `yaml:"date"`
	Slug  string    `yaml:"slug"`
}

func LoadPosts() ([]Post, error) {
	var posts []Post

	files, err := filepath.Glob("content/posts/*.md")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var meta frontmatterData
		body, err := frontmatter.Parse(bytes.NewReader(data), &meta)
		if err != nil {
			println("failed to parse:", file)
			continue
		}

		var buf bytes.Buffer
		err = goldmark.Convert(body, &buf)
		if err != nil {
			continue
		}

		post := Post{
			Title:   meta.Title,
			Date:    meta.Date,
			Slug:    meta.Slug,
			Content: template.HTML(buf.String()),
		}

		posts = append(posts, post)
	}

	return posts, nil
}
