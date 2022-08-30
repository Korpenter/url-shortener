package db

import (
	"fmt"
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

func (r *InMemRepo) GetByID(id string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	v, ok := r.urls[id]
	if !ok {
		return v, fmt.Errorf("unknown id: %s", id)
	}
	return v, nil
}

func (r *InMemRepo) AddLink(long string, id string) {
	r.Lock()
	defer r.Unlock()
	r.urls[id] = long
}
