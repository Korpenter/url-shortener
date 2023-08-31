// Package router provides router for
package router

import (
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/router/handlers"
	"github.com/Mldlr/url-shortener/internal/app/router/middleware"
	"github.com/Mldlr/url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// NewRouter initializes a chi router instance.
func NewRouter(shortener service.ShortenerService, c *config.Config) chi.Router {
	// Initialize new loader to handle batch delete requests.

	r := chi.NewRouter()

	// Define used middlewares for all routes.
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.Decompress)
	r.Use(middleware.Auth{Config: c}.Authenticate)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(chiMiddleware.Compress(5, "application/json", "text/plain"))

	// Define routes.
	r.Mount("/debug", chiMiddleware.Profiler())
	r.Get("/api/user/urls", handlers.APIUserExpand(shortener))
	r.Post("/api/shorten", handlers.APIShorten(shortener))
	r.Post("/api/shorten/batch", handlers.APIShortenBatch(shortener))
	r.Delete("/api/user/urls", handlers.APIDeleteBatch(shortener))
	r.Get("/ping", handlers.Ping(shortener))
	r.Get("/{id}", handlers.Expand(shortener))
	r.Post("/", handlers.Shorten(shortener))
	r.Group(func(r chi.Router) {
		// Define internal route and middleware for it.
		r.Use(middleware.Trusted{Config: c}.TrustCheck)
		r.Get("/api/internal/stats", handlers.APIInternalStats(shortener))
	})
	return r
}
