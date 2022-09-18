package storage

import (
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/config"
	"log"
)

// Repository interface for storage instances
type Repository interface {
	Get(id string) (string, error)
	Add(long, id string) (string, error)
	NewID() (int, error)
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
