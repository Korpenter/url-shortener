package storage

import (
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"sync"
)

type InMemRepo struct {
	urls map[string]string
	sync.RWMutex
}

func NewInMemRepo() *InMemRepo {
	return &InMemRepo{
		urls: make(map[string]string),
	}
}

func (r *InMemRepo) Get(id string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	v, ok := r.urls[id]
	if !ok {
		return v, fmt.Errorf("unknown id: %s", id)
	}
	return v, nil
}

func (r *InMemRepo) Add(longURL string) string {
	r.Lock()
	defer r.Unlock()
	id := encoders.ToBase62(r.NewID())
	r.urls[id] = longURL
	return id
}

func (r *InMemRepo) NewID() int {
	return len(r.urls) + 1
}
