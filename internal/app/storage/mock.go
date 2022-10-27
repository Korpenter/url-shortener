package storage

import (
	"context"
	"fmt"

	"github.com/Mldlr/url-shortener/internal/app/model"
)

type mockRepo struct {
	urlsByShort  map[string]*model.URL
	urlsByUser   map[string][]*model.URL
	existingURLs map[string]*model.URL
	lastID       int
}

// NewMockRepo returns a pointer to a new mock repo instance
func NewMockRepo() *mockRepo {
	mock := mockRepo{
		urlsByShort:  make(map[string]*model.URL),
		urlsByUser:   make(map[string][]*model.URL),
		existingURLs: make(map[string]*model.URL),
	}
	url1 := &model.URL{ShortURL: "1", LongURL: "https://github.com/Mldlr/url-shortener/internal/app/utils/encoders"}
	url2 := &model.URL{ShortURL: "2", LongURL: "https://yandex.ru/"}
	mock.urlsByShort["1"] = url1
	mock.urlsByShort["2"] = url2
	mock.lastID = 2
	mock.urlsByUser["KS097f1lS&F"] = []*model.URL{url1, url2}
	return &mock
}

func (r *mockRepo) Get(short string, ctx context.Context) (string, error) {
	v, ok := r.urlsByShort[short]
	if !ok {
		return "", fmt.Errorf("invalid id: %s", short)
	}
	return v.LongURL, nil
}

func (r *mockRepo) Add(url *model.URL, ctx context.Context) (bool, error) {
	if v, k := r.existingURLs[url.LongURL]; k {
		url.ShortURL = v.ShortURL
		return true, nil
	}
	r.urlsByShort[url.ShortURL] = url
	r.urlsByUser[url.UserID] = append(r.urlsByUser[url.UserID], url)
	r.existingURLs[url.LongURL] = url
	return false, nil
}

func (r *mockRepo) AddBatch(urls map[string]*model.URL, ctx context.Context) (bool, error) {
	var duplicates bool
	for _, v := range urls {
		if i, k := r.existingURLs[v.LongURL]; k {
			duplicates = true
			v.ShortURL = i.ShortURL
			continue
		}
		r.existingURLs[v.LongURL] = v
		r.urlsByShort[v.ShortURL] = v
		r.urlsByUser[v.UserID] = append(r.urlsByUser[v.UserID], v)
	}
	return duplicates, nil
}

func (r *mockRepo) NewID() (int, error) {
	r.lastID++
	return r.lastID, nil
}

func (r *mockRepo) GetByUser(userID string, ctx context.Context) ([]*model.URL, error) {
	s := []*model.URL{}
	s = append(s, r.urlsByUser[userID]...)
	if len(s) == 0 {
		return nil, fmt.Errorf("no urls found for user")
	}
	return s, nil
}

func (r *mockRepo) Ping(context.Context) error {
	return nil
}

func (r *mockRepo) DeleteRepo(context.Context) error {
	r.urlsByShort = make(map[string]*model.URL)
	r.urlsByUser = make(map[string][]*model.URL)
	return nil
}
