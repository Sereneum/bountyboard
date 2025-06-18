package handlers

import (
	"bountyboard/internal/domain/task"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
)

type TaskHandler struct {
	service task.Service
	tmpl    *template.Template
}

func NewTaskHandler(service task.Service, tmpl *template.Template) *TaskHandler {
	return &TaskHandler{service: service, tmpl: tmpl}
}

// List возвращает JSON-массив задач
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	// slog.Info("request", slog.String("url", r.URL.Path))
	userID := "demo" // TODO: заменить на реальный userID из аутентификации

	tasks, err := h.service.ListTasks(userID)
	if err != nil {
		http.Error(w, "failed to load tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "failed to encode tasks", http.StatusInternalServerError)
	}
}

// createTaskRequest структура для декодирования JSON-запроса
type createTaskRequest struct {
	UserID       string `json:"user_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	BountyAmount int    `json:"bounty_amount"`
}

// Add создает новую задачу
func (h *TaskHandler) Add(w http.ResponseWriter, r *http.Request) {
	// slog.Info("request", slog.String("url", r.URL.Path))
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Info("invalid request body", slog.String("err", err.Error()))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.Title == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateTask(req.UserID, req.Title, req.Description, req.BountyAmount); err != nil {
		slog.Info("failed to create task", slog.String("err", err.Error()))
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(`{"status":"created"}`))
}
