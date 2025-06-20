package app

import (
	"bountyboard/internal/adapter/http/renderer"
	"bountyboard/internal/domain/auth"
	"bountyboard/internal/domain/task"
	"net/http"
)

type App struct {
	Router http.Handler
}

type Config struct {
	Repo     task.Repository
	Cache    task.Cache
	Renderer renderer.Renderer
	Auth     auth.Service
}

func Setup(cfg Config) (*App, error) {
	taskService := task.New(cfg.Repo, cfg.Cache)
	router := InitRouter(taskService, cfg.Auth, cfg.Renderer)

	return &App{
		Router: router,
	}, nil
}
