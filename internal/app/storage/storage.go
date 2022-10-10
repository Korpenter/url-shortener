package storage

import (
	"fmt"
	"log"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/model"
)

// Repository interface for storage instances
type Repository interface {
	Get(id string) (string, error)
	GetByUser(userID string) ([]*model.URL, error)
	Add(long, short, userID string) (string, error)
	NewID() (int, error)
	Ping() error
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
		err = r.Ping()
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
		return r
	}
	return NewInMemRepo()
}
