package handlers

import (
	"bountyboard/internal/domain/task"
	"html/template"
)

type Handlers struct {
	Task   *TaskHandler
	Static *StaticHandler
}

type baseHandler struct {
	service task.Service
	tmpl    *template.Template
}

func NewHandlers(service task.Service, tmpl *template.Template) *Handlers {
	base := &baseHandler{service: service, tmpl: tmpl}

	return &Handlers{
		Task:   &TaskHandler{baseHandler: base},
		Static: &StaticHandler{baseHandler: base},
	}
}
