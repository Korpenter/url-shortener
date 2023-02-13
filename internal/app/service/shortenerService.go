package service

import (
	"context"

	"github.com/Mldlr/url-shortener/internal/app/models"
)

// ShortenerService is an interface for handling internal logic of the app
type ShortenerService interface {
	Shorten(ctx context.Context, url *models.URL) (*models.URL, error)
	Expand(ctx context.Context, id string) (*models.URL, error)
	ExpandUser(ctx context.Context, userID string) ([]*models.URL, error)
	APIDeleteBatch(urlIDs []string, userID string)
	Ping(ctx context.Context) error
	ShortenBatch(ctx context.Context, userID string, urls []*models.URL) ([]*models.URL, error)
	Stats(ctx context.Context) (*models.Stats, error)
	BuildURL(url string) string
}
