package app

import (
	"bountyboard/internal/adapter/http/handlers"
	authMiddleware "bountyboard/internal/adapter/http/middleware"
	"bountyboard/internal/adapter/http/renderer"
	"bountyboard/internal/domain/auth"
	"bountyboard/internal/domain/task"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func InitRouter(
	taskService task.Service,
	authService auth.Service,
	renderer renderer.Renderer,
) http.Handler {
	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)
	staticHandler := handlers.NewStaticHandler(taskService, renderer)

	r := chi.NewRouter()

	// === Middleware ===
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// === Public static files ===
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	// === SSR route (HTML-шаблоны) ===
	r.Get("/", staticHandler.Main)
	r.Get("/profile", staticHandler.Profile)

	// === Auth (если есть) ===
	r.Post("/login", authHandler.Login) // POST /login
	// TODO - register
	r.Post("/register", authHandler.Register) // POST /login

	// === API v1 ===
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/tasks", func(r chi.Router) {
			r.Use(authMiddleware.Auth(authService))
			// r.Use(AuthMiddleware) // подключи при необходимости
			r.Get("/", taskHandler.List) // GET    /api/v1/tasks
			r.Post("/", taskHandler.Add) // POST   /api/v1/tasks

			/* TODO - put/delete
			r.Route("/{id}", func(r chi.Router) {
				r.Put("/", taskHandler.Update)    // PUT    /api/v1/tasks/{id}
				r.Delete("/", taskHandler.Delete) // DELETE /api/v1/tasks/{id}
			})

			*/
		})
	})

	return r
}
