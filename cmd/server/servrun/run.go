package servrun

import (
	"bountyboard/internal/adapter/http/renderer"
	"bountyboard/internal/app/factory"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	cachePath = "task-cache.gob"
)

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO renderer -> jet
	r := renderer.NewRenderer()
	layout := filepath.Join("web", "templates", "layout.html")
	if err := r.Add(layout, filepath.Join("web", "templates", "tasks.html"), "tasks"); err != nil {
		return err
	}
	if err := r.Add(layout, filepath.Join("web", "templates", "profile.html"), "profile"); err != nil {
		return err
	}
	//
	//tmpl, err := template.ParseFiles(
	//	filepath.Join("web", "templates", "layout.html"),
	//	filepath.Join("web", "templates", "tasks.html"),
	//	filepath.Join("web", "templates", "profile.html"),
	//)
	//if err != nil {
	//	return err
	//}

	f := factory.NewFactory(ctx, "task-cache.gob", r)

	app, cache, err := f.BuildApp()
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:         ":8080",
		Handler:      app.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("Starting server", slog.String("addr", server.Addr))
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err = server.Shutdown(ctxShutdown); err != nil {
		slog.Error("server shutdown error", slog.String("err", err.Error()))
	} else {
		slog.Info("HTTP server shutdown completed")
	}

	if err = cache.SaveToFile(cachePath); err != nil {
		slog.Error("failed to save task cache", slog.String("err", err.Error()))
	} else {
		slog.Info("task cache saved")
	}

	return nil
}
