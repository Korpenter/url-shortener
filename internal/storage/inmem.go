package storage

import (
	"fmt"
	"github.com/Mldlr/url-shortener/internal/utils/encode"
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

func (r *InMemRepo) Add(long string) string {
	r.Lock()
	defer r.Unlock()
	id := encode.ToBase62(r.NewID())
	r.urls[id] = long
	return id
}

func (r *InMemRepo) NewID() int {
	return len(r.urls)
}
