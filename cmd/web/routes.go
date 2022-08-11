package main

import (
	"net/http"

	"example.com/hello-world/pkg/config"
	"example.com/hello-world/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(a *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// Here we use middlewares
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// Here we use handlers
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	return mux
}