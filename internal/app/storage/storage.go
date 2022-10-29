package storage

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"log"
	"time"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/model"
)

// Repository interface for storage instances
type Repository interface {
	Get(ctx context.Context, id string) (*model.URL, error)
	GetByUser(ctx context.Context, userID string) ([]*model.URL, error)
	Add(ctx context.Context, url *model.URL) (bool, error)
	AddBatch(ctx context.Context, urls map[string]*model.URL) (bool, error)
	NewID(url string) (string, error)
	Ping(ctx context.Context) error
	DeleteRepo(ctx context.Context) error
	DeleteURLs(deleteURLs []*model.DeleteURLItem) (int, error)
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
			r.updateFile()
		})
		s.StartAsync()
		return r
	}
	return NewInMemRepo()
}
