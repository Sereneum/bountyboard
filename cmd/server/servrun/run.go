package servrun

import (
	"bountyboard/internal/app/factory"
	"bountyboard/internal/domain/task"
	"context"
	"errors"
	"html/template"
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

	tmpl, err := template.ParseFiles(
		filepath.Join("web", "templates", "layout.html"),
		filepath.Join("web", "templates", "tasks.html"),
	)
	if err != nil {
		return err
	}

	f := factory.NewFactory(ctx, "task-cache.gob", tmpl)

	app, err := f.BuildApp()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.Handlers.Static.Main)
	mux.HandleFunc("/list", app.Handlers.Task.List)
	mux.HandleFunc("/add", app.Handlers.Task.Add)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
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

	if err := server.Shutdown(ctxShutdown); err != nil {
		slog.Error("server shutdown error", slog.String("err", err.Error()))
	} else {
		slog.Info("HTTP server shutdown completed")
	}

	if fc, ok := app.TaskService.Cache().(task.FileCache); ok {
		if err = fc.SaveToFile(cachePath); err != nil {
			slog.Error("failed to save task cache", slog.String("err", err.Error()))
		} else {
			slog.Info("task cache saved")
		}
	}

	return nil
}
