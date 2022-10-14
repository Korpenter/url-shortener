package storage

import (
	"fmt"
	"sync"

	"github.com/Mldlr/url-shortener/internal/app/model"
)

// InMemRepo is an in-memory url storage
type InMemRepo struct {
	existingURLs map[string]*model.URL
	urlsByShort  map[string]*model.URL
	urlsByUser   map[string][]*model.URL
	lastID       int
	sync.RWMutex
}

// NewInMemRepo returns a pointer to a new repo instance
func NewInMemRepo() *InMemRepo {
	return &InMemRepo{
		urlsByShort:  make(map[string]*model.URL),
		urlsByUser:   make(map[string][]*model.URL),
		existingURLs: make(map[string]*model.URL),
	}
}

// Get returns original link by id or an error if id is not present
func (r *InMemRepo) Get(short string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	v, ok := r.urlsByShort[short]
	if !ok {
		return "", fmt.Errorf("invalid id: %s", short)
	}
	return v.LongURL, nil
}

// Add adds a link to db and returns assigned id
func (r *InMemRepo) Add(url *model.URL) (bool, error) {
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

func (r *InMemRepo) AddBatch(urls map[string]*model.URL) (bool, error) {
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
func (r *InMemRepo) NewID() (int, error) {
	r.Lock()
	defer r.Unlock()
	r.lastID++
	return r.lastID, nil
}

func (r *InMemRepo) GetByUser(userID string) ([]*model.URL, error) {
	r.RLock()
	defer r.RUnlock()
	urls := make([]*model.URL, 0)
	urls = append(urls, r.urlsByUser[userID]...)
	if len(urls) == 0 {
		return nil, nil
	}
	return urls, nil
}

func (r *InMemRepo) Ping() error {
	return nil
}

func (r *InMemRepo) DeleteRepo() error {
	r.Lock()
	defer r.Unlock()
	r.urlsByShort = make(map[string]*model.URL)
	r.urlsByUser = make(map[string][]*model.URL)
	return nil
}
