package loader

import (
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"time"
)

func NewDeleteLoader(repo storage.Repository) *UserLoader {
	deleteLoaderCfg := UserLoaderConfig{
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
	return NewUserLoader(deleteLoaderCfg)
}
