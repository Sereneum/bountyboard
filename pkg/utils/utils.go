package utils

import (
	"bytes"
	"github.com/yuin/goldmark"
	"html/template"
)

func MdToHTML(md string) template.HTML {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return template.HTML("<p><em>Error parsing description</em></p>")
	}
	return template.HTML(buf.String())
}
