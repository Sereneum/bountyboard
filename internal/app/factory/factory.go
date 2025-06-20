package factory

import (
	"bountyboard/internal/adapter/cache/memory"
	"bountyboard/internal/adapter/http/renderer"
	inmemory "bountyboard/internal/adapter/storage/in-memory"
	"bountyboard/internal/app"
	"bountyboard/internal/domain/auth"
	"bountyboard/internal/domain/task"
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type Factory struct {
	ctx       context.Context
	cachePath string // TODO - cfg
	r         renderer.Renderer
}

func NewFactory(ctx context.Context, cachePath string, r renderer.Renderer) *Factory {
	return &Factory{
		ctx:       ctx,
		cachePath: cachePath,
		r:         r,
	}
}

func (f *Factory) BuildApp() (*app.App, task.FileCache, error) {
	// 1. Репозиторий
	repo := inmemory.NewRepo()

	// 2. Кэш
	cache := memory.NewTaskCache(7 * 24 * time.Hour)
	go cache.Janitor(f.ctx, 10*time.Minute)
	go cache.RunAutoSave(f.ctx, f.cachePath, time.Hour)

	if err := cache.LoadFromFile(f.cachePath); err != nil {
		if !os.IsNotExist(err) { // Игнорируем если файла нет
			return nil, nil, fmt.Errorf("failed to load cache: %w", err)
		}
		slog.Warn("cache file not found, starting with empty cache")
	}

	authService := auth.New("secret") // TODO os.Getenv("JWT_SECRET")

	// 3. Настраиваем и возвращаем App
	appCfg := app.Config{
		Repo:     repo,
		Cache:    cache,
		Renderer: f.r,
		Auth:     authService,
	}

	a, err := app.Setup(appCfg)
	if err != nil {
		return nil, nil, err
	}

	return a, cache, nil
}
