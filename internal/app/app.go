package app

import (
	"bountyboard/internal/adapter/http/handlers"
	"bountyboard/internal/domain/task"
	"html/template"
)

type App struct {
	TaskService task.Service
	Handler     *handlers.TaskHandler
	Templates   *template.Template
}

func Setup(repo task.Repository, cache task.Cache, tmpl *template.Template) (*App, error) {
	// Создаем сервис из repo и cache
	service := task.New(repo, cache)

	// Создаем хендлер, передаем сервис и шаблоны
	handler := handlers.NewHandler(service, tmpl)

	return &App{
		TaskService: service,
		Handler:     handler,
		Templates:   tmpl,
	}, nil
}

/*
func Setup(ctx context.Context, cachePath string, tmpl *template.Template) (*App, error) {
	// репозиторий
	repo := inmemeory.NewRepo()
	slog.Info("init", slog.String("entry", "in-memory repo"))

	// кэш
	cache := memory.NewTaskCache(7 * 24 * time.Hour) // concrete type implementing FileCache
	go cache.Janitor(ctx, 10*time.Minute)
	go cache.RunAutoSave(ctx, cachePath, time.Hour)
	if err := cache.LoadFromFile(cachePath); err != nil {
		slog.Error("task cache error", slog.String("err", err.Error()))
	}

	// сервис
	service := task.New(repo, cache)
	slog.Info("init", slog.String("entry", "task service"))

	// хэндлер
	handler := handlers.NewHandler(service, tmpl)
	slog.Info("init", slog.String("entry", "task handler"))

	return &App{
		TaskService: service,
		Handler:     handler,
	}, nil
}
*/
