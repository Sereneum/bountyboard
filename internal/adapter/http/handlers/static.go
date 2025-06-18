package handlers

import (
	"bountyboard/internal/domain/task"
	"log/slog"
	"net/http"
)

type StaticHandler struct {
	*baseHandler
}

func (h *StaticHandler) Main(w http.ResponseWriter, r *http.Request) {
	userID := "demo" // TODO: заменить на реальный userID из аутентификации
	tasks, err := h.service.ListTasks(userID)
	if err != nil {
		slog.Error("failed to list tasks", slog.String("err", err.Error()))
		http.Error(w, "failed to load tasks", http.StatusInternalServerError)
		return
	}

	data := struct {
		Tasks []*task.Task
	}{
		Tasks: tasks,
	}

	if err = h.tmpl.Execute(w, data); err != nil {
		slog.Error("template execute error", slog.String("err", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
