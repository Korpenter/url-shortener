package storage

import (
	"fmt"
)

type MockRepo struct {
	urls map[string]string
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
	v, ok := r.urls[id]
	if !ok {
		return v, fmt.Errorf("invalid id: %s", id)
	}
	return v, nil
}

func (r *MockRepo) Add(longURL, id string) (string, error) {
	r.urls[id] = longURL
	return id, nil
}

func (r *MockRepo) NewID() (int, error) {
	return len(r.urls) + 1, nil
}
