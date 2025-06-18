package app

import (
	"bountyboard/internal/domain/auth"
	"bountyboard/internal/domain/task"
	"html/template"
	"net/http"
)

type App struct {
	Router http.Handler
}

type Config struct {
	Repo      task.Repository
	Cache     task.Cache
	Templates *template.Template
	Auth      auth.Service
}

func Setup(cfg Config) (*App, error) {
	taskService := task.New(cfg.Repo, cfg.Cache)
	router := InitRouter(taskService, cfg.Auth, cfg.Templates)

	return &App{
		Router: router,
	}, nil
}
