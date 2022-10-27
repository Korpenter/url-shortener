package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/model"
)

// Repository interface for storage instances
type Repository interface {
	Get(id string, ctx context.Context) (string, error)
	GetByUser(userID string, ctx context.Context) ([]*model.URL, error)
	Add(url *model.URL, ctx context.Context) (bool, error)
	AddBatch(urls map[string]*model.URL, ctx context.Context) (bool, error)
	NewID() (int, error)
	Ping(ctx context.Context) error
	DeleteRepo(ctx context.Context) error
}

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
		err = r.NumberOfURLs()
		if err != nil {
			log.Fatal(fmt.Errorf("error getting number of records : %v", err))
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
		return r
	}
	return NewInMemRepo()
}
