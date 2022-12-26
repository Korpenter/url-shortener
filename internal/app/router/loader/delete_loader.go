package loader

import (
	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"time"
)

func NewDeleteLoader(repo storage.Repository) *UserLoader {
	deleteLoaderCfg := UserLoaderConfig{
		MaxBatch: 200,
		Wait:     5 * time.Second,
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
