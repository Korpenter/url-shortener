package storage

import (
	"fmt"
	"sync"
)

// InMemRepo is an in-memory url storage
type InMemRepo struct {
	urls map[string]string
	sync.RWMutex
}

// NewInMemRepo returns a pointer to a new repo instance
func NewInMemRepo() *InMemRepo {
	return &InMemRepo{
		urls: make(map[string]string),
	}
}

// Get returns original link by id or an error if id is not present
func (r *InMemRepo) Get(id string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	v, ok := r.urls[id]
	if !ok {
		return v, fmt.Errorf("invalid id: %s", id)
	}
	return v, nil
}

// Add adds a link to db and returns assigned id
func (r *InMemRepo) Add(longURL, id string) string {
	r.Lock()
	defer r.Unlock()
	r.urls[id] = longURL
	return id
}

// NewID returns a number to encode as an id
func (r *InMemRepo) NewID() int {
	r.Lock()
	defer r.Unlock()
	return len(r.urls) + 1
}
