package storage

import (
	"fmt"
	"log"

	"github.com/Mldlr/url-shortener/internal/app/config"
)

// Repository interface for storage instances
type Repository interface {
	Get(id string) (string, error)
	Add(long, id string) (string, error)
	NewID() (int, error)
}

// url represents url record
type url struct {
	ID      string `json:"id"`
	LongURL string `json:"long_url"`
}

func New(c *config.Config) Repository {
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
