package renderer

import (
	"html/template"
	"os"
)

type Renderer interface {
	Add(string, string, string) error
	Get(string) (*template.Template, error)
}
type renderer struct {
	tmpl map[string]*template.Template
}

func NewRenderer() Renderer {
	return &renderer{
		tmpl: make(map[string]*template.Template),
	}
}

func (r *renderer) Add(layout, content, name string) error {
	tmpl, err := template.ParseFiles(layout, content)
	if err != nil {
		return nil
	}

	r.tmpl[name] = tmpl

	return nil
}

func (r *renderer) Get(name string) (*template.Template, error) {
	t, found := r.tmpl[name]
	if !found {
		return nil, os.ErrNotExist
	}
	return t, nil
}

//tmpl, err := template.ParseFiles(
//filepath.Join("web", "templates", "layout.html"),
//filepath.Join("web", "templates", "tasks.html"),
//filepath.Join("web", "templates", "profile.html"),
//)
