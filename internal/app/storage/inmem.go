package storage

import (
	"context"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"sync"

	"github.com/Mldlr/url-shortener/internal/app/models"
)

// InMemRepo is an in-memory url storage
type InMemRepo struct {
	existingURLs map[string]*models.URL
	urlsByShort  map[string]*models.URL
	urlsByUser   map[string][]*models.URL
	sync.RWMutex
}

// NewInMemRepo returns a pointer to a new repo instance
func NewInMemRepo() *InMemRepo {
	return &InMemRepo{
		urlsByShort:  make(map[string]*models.URL),
		urlsByUser:   make(map[string][]*models.URL),
		existingURLs: make(map[string]*models.URL),
	}
}

// Get returns original link by id or an error if id is not present
func (r *InMemRepo) Get(ctx context.Context, id string) (*models.URL, error) {
	r.RLock()
	defer r.RUnlock()
	url, ok := r.urlsByShort[id]
	if !ok {
		return nil, fmt.Errorf("invalid id: %s", id)
	}
	return url, nil
}

// Add adds a link to db and returns assigned id
func (r *InMemRepo) Add(ctx context.Context, url *models.URL) (bool, error) {
	r.Lock()
	defer r.Unlock()
	if v, k := r.existingURLs[url.LongURL]; k {
		url.ShortURL = v.ShortURL
		return true, nil
	}
	r.urlsByShort[url.ShortURL] = url
	r.urlsByUser[url.UserID] = append(r.urlsByUser[url.UserID], url)
	r.existingURLs[url.LongURL] = url
	return false, nil
}

func (r *InMemRepo) AddBatch(ctx context.Context, urls map[string]*models.URL) (bool, error) {
	r.Lock()
	defer r.Unlock()
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
func (r *InMemRepo) NewID(url string) (string, error) {
	return encoders.ToRBase62(url), nil
}

func (r *InMemRepo) GetByUser(ctx context.Context, userID string) ([]*models.URL, error) {
	r.RLock()
	defer r.RUnlock()
	urls := make([]*models.URL, 0)
	urls = append(urls, r.urlsByUser[userID]...)
	if len(urls) == 0 {
		return nil, nil
	}
	return urls, nil
}

func (r *InMemRepo) DeleteURLs(deleteURLs []*models.DeleteURLItem) (int, error) {
	r.Lock()
	defer r.Unlock()
	var n int
	for _, v := range deleteURLs {
		if _, ok := r.urlsByShort[v.ShortURL]; ok && r.urlsByShort[v.ShortURL].UserID == v.UserID {
			r.urlsByShort[v.ShortURL].Deleted = true
			n++
		}
	}
	return n, nil
}

func (r *InMemRepo) Ping(context.Context) error {
	return nil
}

func (r *InMemRepo) DeleteRepo(context.Context) error {
	r.Lock()
	defer r.Unlock()
	r.urlsByShort = make(map[string]*models.URL)
	r.urlsByUser = make(map[string][]*models.URL)
	return nil
}
