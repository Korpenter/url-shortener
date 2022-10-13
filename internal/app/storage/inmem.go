package storage

import (
	"fmt"
	"sync"

	"github.com/Mldlr/url-shortener/internal/app/model"
)

// InMemRepo is an in-memory url storage
type InMemRepo struct {
	urlsByShort map[string]*model.URL
	urlsByUser  map[string][]*model.URL
	lastID      int
	sync.RWMutex
}

// NewInMemRepo returns a pointer to a new repo instance
func NewInMemRepo() *InMemRepo {
	return &InMemRepo{
		urlsByShort: make(map[string]*model.URL),
		urlsByUser:  make(map[string][]*model.URL),
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
func (r *InMemRepo) Add(url *model.URL) error {
	r.Lock()
	defer r.Unlock()
	r.urlsByShort[url.ShortURL] = url
	r.urlsByUser[url.UserID] = append(r.urlsByUser[url.UserID], url)
	return nil
}

func (r *InMemRepo) AddBatch(urls map[string]*model.URL) error {
	r.Lock()
	defer r.Unlock()
	for _, v := range urls {
		r.urlsByShort[v.ShortURL] = v
		r.urlsByUser[v.UserID] = append(r.urlsByUser[v.UserID], v)
	}
	return nil
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
