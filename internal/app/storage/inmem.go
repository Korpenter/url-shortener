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
func (r *InMemRepo) Add(longURL, short, userID string) (string, error) {
	r.Lock()
	defer r.Unlock()
	url := &model.URL{ShortURL: short, LongURL: longURL}
	r.urlsByShort[short] = url
	r.urlsByUser[userID] = append(r.urlsByUser[userID], url)
	return short, nil
}

// NewID returns a number to encode as an id
func (r *InMemRepo) NewID() (int, error) {
	r.Lock()
	defer r.Unlock()
	return len(r.urlsByShort) + 1, nil
}

func (r *InMemRepo) GetByUser(userID string) ([]*model.URL, error) {
	r.RLock()
	defer r.RUnlock()
	s := []*model.URL{}
	s = append(s, r.urlsByUser[userID]...)
	if len(s) == 0 {
		return nil, nil
	}
	return s, nil
}
