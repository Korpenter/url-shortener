// Package storage provides storage implementations for url-shortener
package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/models"
)

// Repository is an interface for storage instances
type Repository interface {
	Get(ctx context.Context, id string) (*models.URL, error)
	GetByUser(ctx context.Context, userID string) ([]*models.URL, error)
	Add(ctx context.Context, url *models.URL) (bool, error)
	AddBatch(ctx context.Context, urls map[string]*models.URL) (bool, error)
	NewID(url string) (string, error)
	Ping(ctx context.Context) error
	DeleteRepo(ctx context.Context) error
	DeleteURLs(deleteURLs []*models.DeleteURLItem) (int, error)
	Stats(ctx context.Context) (*models.Stats, error)
	Close() error
}

// New initializes a new Repository instance to use as a storage.
func New(c *config.Config) Repository {
	if c.PostgresURL != "" {
		r, err := NewPostgresRepo(c.PostgresURL)
		if err != nil {
			log.Fatal(fmt.Errorf("error initiating postgres connection : %v", err))
		}
		err = r.NewTableURLs()
		if err != nil {
			log.Fatal(fmt.Errorf("error creating urls table  : %v", err))
		}
		err = r.Ping(context.Background())
		if err != nil {
			log.Fatal(fmt.Errorf("error pinging db : %v", err))
		}
		return r
	}
	if c.FileStorage != "" {
		r, err := NewFileRepo(c.FileStorage)
		if err != nil {
			log.Fatal(fmt.Errorf("error initiating file storage : %v", err))
		}
		err = r.Load()
		if err != nil {
			log.Fatal(fmt.Errorf("error loading json data from file : %v", err))
		}
		s := gocron.NewScheduler(time.UTC)
		s.Every(1).Minutes().Do(func() {
			r.update()
		})
		s.StartAsync()
		return r
	}
	return NewInMemRepo()
}
