package router

import (
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/router/handlers"
	"github.com/Mldlr/url-shortener/internal/app/router/middleware"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// NewRouter returns a chi router instance.
func NewRouter(repo storage.Repository, c *config.Config) chi.Router {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.Decompress)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(chiMiddleware.Compress(5, "application/json", "text/plain"))
	r.Post("/api/shorten", handlers.APIShorten(repo, c))
	r.Get("/{id}", handlers.Expand(repo))
	r.Post("/", handlers.Shorten(repo, c))
	return r
}
