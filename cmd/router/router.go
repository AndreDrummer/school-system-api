package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/AndreDrummer/school-system-api/cmd/router/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func handler() http.Handler {
	handler := chi.NewMux()

	handler.Use(middleware.Recoverer)
	handler.Use(middleware.RequestID)
	handler.Use(middleware.Logger)

	handler.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			routes.Students(r)
		})
	})

	return handler
}

func Run() {
	serverHanlder := handler()
	server := &http.Server{
		Addr:         ":8080",
		Handler:      serverHanlder,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

	err := server.ListenAndServe()

	if err != nil {
		slog.Error("Could not initialize server..")
	}
}
