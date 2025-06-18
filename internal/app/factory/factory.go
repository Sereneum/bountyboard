package factory

import (
	"bountyboard/internal/adapter/cache/memory"
	inmemory "bountyboard/internal/adapter/storage/in-memory"
	"bountyboard/internal/app"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"time"
)

type Factory struct {
	ctx       context.Context
	cachePath string // TODO - cfg
	tmpl      *template.Template
}

func NewFactory(ctx context.Context, cachePath string, tmpl *template.Template) *Factory {
	return &Factory{
		ctx:       ctx,
		cachePath: cachePath,
		tmpl:      tmpl,
	}
}

func (f *Factory) BuildApp() (*app.App, error) {
	// 1. Репозиторий
	repo := inmemory.NewRepo()

	// 2. Кэш
	cache := memory.NewTaskCache(7 * 24 * time.Hour)
	go cache.Janitor(f.ctx, 10*time.Minute)
	go cache.RunAutoSave(f.ctx, f.cachePath, time.Hour)

	if err := cache.LoadFromFile(f.cachePath); err != nil {
		if !os.IsNotExist(err) { // Игнорируем если файла нет
			return nil, fmt.Errorf("failed to load cache: %w", err)
		}
		slog.Warn("cache file not found, starting with empty cache")
	}

	// 3. Настраиваем и возвращаем App
	return app.Setup(repo, cache, f.tmpl)
}
