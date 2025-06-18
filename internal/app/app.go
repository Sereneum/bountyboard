package app

import (
	"bountyboard/internal/adapter/http/handlers"
	"bountyboard/internal/domain/task"
	"html/template"
)

type App struct {
	TaskService task.Service
	Handlers    *handlers.Handlers
}

func Setup(repo task.Repository, cache task.Cache, tmpl *template.Template) (*App, error) {
	// Создаем сервис из repo и cache
	service := task.New(repo, cache)

	// Создаем хендлер, передаем сервис и шаблоны
	h := handlers.NewHandlers(service, tmpl)

	return &App{
		TaskService: service,
		Handlers:    h,
	}, nil
}
