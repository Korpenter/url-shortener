package router

import (
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/Mldlr/url-shortener/internal/app/router/handlers"
	"github.com/Mldlr/url-shortener/internal/app/router/loader"
	"github.com/Mldlr/url-shortener/internal/app/router/middleware"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"time"
	//"golang.org/x/sync/errgroup"
)

// NewRouter returns a chi router instance.
func NewRouter(repo storage.Repository, c *config.Config) chi.Router {

	deleteLoaderCfg := loader.UserLoaderConfig{
		MaxBatch: 200,
		Wait:     5 * time.Second,
		Fetch: func(keys []*model.DeleteURLItem) ([]int, []error) {
			n, err := repo.DeleteURLs(keys)
			if err != nil {
				return []int{n}, []error{err}
			}
			return []int{n}, nil
		},
	}
	deleteLoader := loader.NewUserLoader(deleteLoaderCfg)

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
