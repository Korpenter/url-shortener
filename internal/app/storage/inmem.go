package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"

	"github.com/Mldlr/url-shortener/internal/app/models"
)

// InMemRepo is an in-memory url storage
type InMemRepo struct {
	// existingURLs maps original URLs to their corresponding URL models.
	existingURLs map[string]*models.URL
	// urlsByShort maps short URLs to their corresponding URL models.
	urlsByShort map[string]*models.URL
	// urlsByUser maps user IDs to their corresponding URL models.
	urlsByUser map[string][]*models.URL
	// RWMutex synchronizes access to the FileRepo.
	sync.RWMutex
}

// NewInMemRepo initializes new in-memory storage
func NewInMemRepo() *InMemRepo {
	return &InMemRepo{
		urlsByShort:  make(map[string]*models.URL),
		urlsByUser:   make(map[string][]*models.URL),
		existingURLs: make(map[string]*models.URL),
	}
}

// Get returns original link by ID or an error if id is not present
func (r *InMemRepo) Get(ctx context.Context, id string) (*models.URL, error) {
	r.RLock()
	defer r.RUnlock()
	url, ok := r.urlsByShort[id]
	if !ok {
		return nil, fmt.Errorf("invalid id: %s", id)
	}
	return url, nil
}

// Add adds a link to storage.
func (r *InMemRepo) Add(ctx context.Context, url *models.URL) (bool, error) {
	r.Lock()
	defer r.Unlock()
	// Check for url in map and return if it already exists.
	if v, k := r.existingURLs[url.LongURL]; k {
		url.ShortURL = v.ShortURL
		return true, nil
	}
	// Otherwise add url to maps.
	r.urlsByShort[url.ShortURL] = url
	r.urlsByUser[url.UserID] = append(r.urlsByUser[url.UserID], url)
	r.existingURLs[url.LongURL] = url
	return false, nil
}

// AddBatch adds multiple URLs to storage.
func (r *InMemRepo) AddBatch(ctx context.Context, urls []*models.URL) (bool, error) {
	r.Lock()
	defer r.Unlock()
	var duplicates bool
	// For each url check if url is in map and add it otherwise.
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
func (r *InMemRepo) NewID(url string) (string, error) {
	return encoders.ToRBase62(url), nil
}

// GetByUser finds URLs created by user.
func (r *InMemRepo) GetByUser(ctx context.Context, userID string) ([]*models.URL, error) {
	r.RLock()
	defer r.RUnlock()
	urls := make([]*models.URL, 0)
	// Get all urls for user.
	urls = append(urls, r.urlsByUser[userID]...)
	if len(urls) == 0 {
		return nil, nil
	}
	return urls, nil
}

// DeleteURLs delete urls from maps.
func (r *InMemRepo) DeleteURLs(deleteURLs []*models.DeleteURLItem) (int, error) {
	r.Lock()
	defer r.Unlock()
	var n int
	// For each of the urls check if the user created this url and delete it if confirmed.
	for _, v := range deleteURLs {
		if _, ok := r.urlsByShort[v.ShortURL]; ok && r.urlsByShort[v.ShortURL].UserID == v.UserID {
			r.urlsByShort[v.ShortURL].Deleted = true
			n++
		}
	}
	return n, nil
}

// Ping is redundant for in-memory storage.
func (r *InMemRepo) Ping(context.Context) error {
	return nil
}

// Stats gets count of urls and registered users
func (r *InMemRepo) Stats(ctx context.Context) (*models.Stats, error) {
	r.RLock()
	defer r.RUnlock()
	var stats models.Stats
	stats.URLCount = len(r.existingURLs)
	stats.UserCount = len(r.urlsByUser)
	return &stats, nil
}

// DeleteRepo deletes repository data.
func (r *InMemRepo) DeleteRepo(context.Context) error {
	r.Lock()
	defer r.Unlock()
	// Reallocate maps.
	r.urlsByShort = make(map[string]*models.URL)
	r.urlsByUser = make(map[string][]*models.URL)
	return nil
}

// Close in not implemented for inmem
func (r *InMemRepo) Close() error {
	return nil
}
