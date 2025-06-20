package handlers

import (
	"bountyboard/internal/adapter/http/renderer"
	"bountyboard/internal/domain/task"
	"log/slog"
	"net/http"
)

//type StaticHandler struct {
//	*baseHandler
//}

type StaticHandler struct {
	taskService task.Service
	renderer    renderer.Renderer
}

func NewStaticHandler(taskService task.Service, renderer renderer.Renderer) *StaticHandler {
	return &StaticHandler{taskService: taskService, renderer: renderer}
}

func tmplError(w http.ResponseWriter, err error) {
	slog.Error("template execute error", slog.String("err", err.Error()))
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

func (h *StaticHandler) Main(w http.ResponseWriter, r *http.Request) {
	userID := "demo" // TODO: заменить на реальный userID из аутентификации
	tasks, err := h.taskService.ListTasks(userID)
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
	//h.tmpl.Execute(w, data)

	//if err = h.tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
	//	slog.Error("template execute error", slog.String("err", err.Error()))
	//	http.Error(w, "internal server error", http.StatusInternalServerError)
	//}
	// TODO - тут точно пора переходить на jet
	t, err := h.renderer.Get("tasks")
	if err != nil {
		tmplError(w, err)
	}

	if t != nil {
		if err = t.Execute(w, data); err != nil {
			tmplError(w, err)
		}
	} else {
		tmplError(w, err)
	}
}

func (h *StaticHandler) Profile(w http.ResponseWriter, r *http.Request) {
	// userID := "demo"

	t, err := h.renderer.Get("profile")
	if err != nil {
		tmplError(w, err)
	}

	if t != nil {
		if err = t.Execute(w, struct{}{}); err != nil {
			tmplError(w, err)
		}
	} else {
		tmplError(w, err)
	}
}
