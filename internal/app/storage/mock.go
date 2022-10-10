package storage

import (
	"fmt"

	"github.com/Mldlr/url-shortener/internal/app/model"
)

type mockRepo struct {
	urlsByShort map[string]*model.URL
	urlsByUser  map[string][]*model.URL
}

// NewMockRepo returns a pointer to a new mock repo instance
func NewMockRepo() *mockRepo {
	mock := mockRepo{
		urlsByShort: make(map[string]*model.URL),
		urlsByUser:  make(map[string][]*model.URL),
	}
	url1 := &model.URL{ShortURL: "1", LongURL: "https://github.com/Mldlr/url-shortener/internal/app/utils/encoders"}
	url2 := &model.URL{ShortURL: "2", LongURL: "https://yandex.ru/"}
	mock.urlsByShort["1"] = url1
	mock.urlsByShort["2"] = url2
	mock.urlsByUser["KS097f1lS&F"] = []*model.URL{url1, url2}
	return &mock
}

func (r *mockRepo) Get(short string) (string, error) {
	v, ok := r.urlsByShort[short]
	if !ok {
		return "", fmt.Errorf("invalid id: %s", short)
	}
	return v.LongURL, nil
}

func (r *mockRepo) Add(longURL, short, userID string) (string, error) {
	url := &model.URL{ShortURL: short, LongURL: longURL}
	r.urlsByShort[short] = url
	r.urlsByUser[userID] = append(r.urlsByUser[userID], url)
	return short, nil
}

func (r *mockRepo) NewID() (int, error) {
	return len(r.urlsByShort) + 1, nil
}

func (r *mockRepo) GetByUser(userID string) ([]*model.URL, error) {
	s := []*model.URL{}
	s = append(s, r.urlsByUser[userID]...)
	if len(s) == 0 {
		return nil, fmt.Errorf("no urls found for user")
	}
	return s, nil
}

func (r *mockRepo) Ping() error {
	return nil
}
