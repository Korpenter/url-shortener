// Package loader provides the loader to handle batch requests.
package loader

import (
	"time"

	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/storage"
)

// NewDeleteLoader creates a new UserLoader instance used to delete URLs in batches.
func NewDeleteLoader(repo storage.Repository) *UserLoader {
	deleteLoaderCfg := UserLoaderConfig{
		// Conditions of starting batch delete.
		MaxBatch: 200,
		Wait:     5 * time.Second,
		// Batch delete function that returns the amount of deleted urls.
		Fetch: func(keys []*models.DeleteURLItem) ([]int, []error) {
			n, err := repo.DeleteURLs(keys)
			if err != nil {
				return []int{n}, []error{err}
			}
			return []int{n}, nil
		},
	}
	return NewUserLoader(deleteLoaderCfg)
}
