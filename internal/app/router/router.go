package router

import (
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/router/handlers"
	"github.com/Mldlr/url-shortener/internal/app/router/loader"
	"github.com/Mldlr/url-shortener/internal/app/router/middleware"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// NewRouter returns a chi router instance.
func NewRouter(repo storage.Repository, c *config.Config) chi.Router {
	deleteLoader := loader.NewDeleteLoader(repo)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.Decompress)
	r.Use(middleware.Auth{Config: c}.Authenticate)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(chiMiddleware.Compress(5, "application/json", "text/plain"))
	r.Get("/api/user/urls", handlers.APIUserExpand(repo, c))
	r.Post("/api/shorten", handlers.APIShorten(repo, c))
	r.Post("/api/shorten/batch", handlers.APIShortenBatch(repo, c))
	r.Delete("/api/user/urls", handlers.APIDeleteBatch(deleteLoader))
	r.Get("/{id}", handlers.Expand(repo))
	r.Get("/ping", handlers.Ping(repo))
	r.Post("/", handlers.Shorten(repo, c))
	return r
}
