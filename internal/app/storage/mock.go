package storage

import (
	"fmt"
	"sync"
)

type MockRepo struct {
	urls map[string]string
	sync.RWMutex
}

// NewMockRepo returns a pointer to a new mock repo instance
func NewMockRepo() *MockRepo {
	mock := MockRepo{
		urls: make(map[string]string),
	}

	mock.urls["1"] = "https://github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	mock.urls["2"] = "https://yandex.ru/"
	return &mock
}

func (r *MockRepo) Get(id string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	v, ok := r.urls[id]
	if !ok {
		return v, fmt.Errorf("invalid id: %s", id)
	}
	return v, nil
}

func (r *MockRepo) Add(longURL, id string) string {
	r.Lock()
	defer r.Unlock()
	r.urls[id] = longURL
	return id
}

func (r *MockRepo) NewID() int {
	r.Lock()
	defer r.Unlock()
	return len(r.urls) + 1
}
