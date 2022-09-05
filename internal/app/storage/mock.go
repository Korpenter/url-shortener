package storage

import (
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"sync"
)

type MockRepo struct {
	urls map[string]string
	sync.RWMutex
}

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
		return v, fmt.Errorf("unknown id: %s", id)
	}
	return v, nil
}

func (r *MockRepo) Add(longURL string) string {
	r.Lock()
	defer r.Unlock()
	id := encoders.ToRBase62(r.NewID())
	r.urls[id] = longURL
	return id
}

func (r *MockRepo) NewID() int {
	return len(r.urls) + 1
}
