package storage

import (
	"fmt"
)

type mockRepo struct {
	urls map[string]string
}

// NewMockRepo returns a pointer to a new mock repo instance
func NewMockRepo() *mockRepo {
	mock := mockRepo{
		urls: make(map[string]string),
	}

	mock.urls["1"] = "https://github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	mock.urls["2"] = "https://yandex.ru/"
	return &mock
}

func (r *mockRepo) Get(id string) (string, error) {
	v, ok := r.urls[id]
	if !ok {
		return v, fmt.Errorf("invalid id: %s", id)
	}
	return v, nil
}

func (r *mockRepo) Add(longURL, id string) (string, error) {
	r.urls[id] = longURL
	return id, nil
}

func (r *mockRepo) NewID() (int, error) {
	return len(r.urls) + 1, nil
}
