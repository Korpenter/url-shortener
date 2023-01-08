package storage

import (
	"context"
	"fmt"

	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
)

type mockRepo struct {
	urlsByShort  map[string]*models.URL
	urlsByUser   map[string][]*models.URL
	existingURLs map[string]*models.URL
}

// NewMockRepo initializes new test in-memory storage with data.
func NewMockRepo() *mockRepo {
	mock := mockRepo{
		urlsByShort:  make(map[string]*models.URL),
		urlsByUser:   make(map[string][]*models.URL),
		existingURLs: make(map[string]*models.URL),
	}
	url1 := &models.URL{ShortURL: "1", LongURL: "https://github.com/Mldlr/url-shortener/internal/app/utils/encoders"}
	url2 := &models.URL{ShortURL: "2", LongURL: "https://yandex.ru/"}
	mock.urlsByShort["1"] = url1
	mock.urlsByShort["2"] = url2
	mock.urlsByUser["KS097f1lS&F"] = []*models.URL{url1, url2}
	return &mock
}

// Get returns original link by ID or an error if id is not present
func (r *mockRepo) Get(ctx context.Context, id string) (*models.URL, error) {
	url, ok := r.urlsByShort[id]
	if !ok {
		return nil, fmt.Errorf("invalid id: %s", id)
	}
	return url, nil
}

// Add adds a link to storage.
func (r *mockRepo) Add(ctx context.Context, url *models.URL) (bool, error) {
	if v, k := r.existingURLs[url.LongURL]; k {
		url.ShortURL = v.ShortURL
		return true, nil
	}
	r.urlsByShort[url.ShortURL] = url
	r.urlsByUser[url.UserID] = append(r.urlsByUser[url.UserID], url)
	r.existingURLs[url.LongURL] = url
	return false, nil
}

// AddBatch adds multiple URLs to storage.
func (r *mockRepo) AddBatch(ctx context.Context, urls map[string]*models.URL) (bool, error) {
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

// NewID returns a number to encode as an id
func (r *mockRepo) NewID(url string) (string, error) {
	return encoders.ToRBase62(url), nil
}

// GetByUser finds URLs created by user.
func (r *mockRepo) GetByUser(ctx context.Context, userID string) ([]*models.URL, error) {
	s := make([]*models.URL, 0)
	s = append(s, r.urlsByUser[userID]...)
	if len(s) == 0 {
		return nil, fmt.Errorf("no urls found for user")
	}
	return s, nil
}

// DeleteURLs delete urls from maps.
func (r *mockRepo) DeleteURLs(deleteURLs []*models.DeleteURLItem) (int, error) {
	var n int
	for _, v := range deleteURLs {
		if _, ok := r.urlsByShort[v.ShortURL]; ok && r.urlsByShort[v.ShortURL].UserID == v.UserID {
			r.urlsByShort[v.ShortURL].Deleted = true
			n++
		}
	}
	return n, nil
}

// Ping is redundant for in-memory storage.
func (r *mockRepo) Ping(context.Context) error {
	return nil
}

// DeleteRepo deletes repository data.
func (r *mockRepo) DeleteRepo(context.Context) error {
	r.urlsByShort = make(map[string]*models.URL)
	r.urlsByUser = make(map[string][]*models.URL)
	return nil
}
