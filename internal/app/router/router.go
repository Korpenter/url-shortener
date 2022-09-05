package router

import (
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/router/handlers"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

func NewRouter(repo storage.Repository, c *config.Config) chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", handlers.Expand(repo))
	r.Post("/", handlers.Shorten(repo, c))
	return r
}
