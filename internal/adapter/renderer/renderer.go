package renderer

// TODO: use in handlers

import (
	"bountyboard/internal/domain/task"
	"html/template"
	"net/http"
)

type Renderer struct {
	tmpl *template.Template
}

func NewRenderer(tmpl *template.Template) *Renderer {
	return &Renderer{tmpl: tmpl}
}

func (r *Renderer) RenderTasks(w http.ResponseWriter, tasks []*task.Task) {
	data := struct {
		Tasks []*task.Task
	}{
		Tasks: tasks,
	}

	if err := r.tmpl.Execute(w, data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}
