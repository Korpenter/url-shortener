package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/Mldlr/url-shortener/internal/app/model"
)

// FileRepo is an in-file url storage
type FileRepo struct {
	file         *os.File
	cacheByShort map[string]*model.URL
	cacheByUser  map[string][]*model.URL
	existingURLs map[string]*model.URL
	encoder      json.Encoder
	lastID       int
	sync.RWMutex
}

// NewFileRepo returns a pointer to a new repo instance
func NewFileRepo(filename string) (*FileRepo, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error openin file : %v", err)
	}
	return &FileRepo{
		file:         file,
		cacheByShort: make(map[string]*model.URL),
		cacheByUser:  make(map[string][]*model.URL),
		existingURLs: make(map[string]*model.URL),
		encoder:      *json.NewEncoder(file),
	}, nil
}

// Load loads stored url records from file
func (r *FileRepo) Load() error {
	decoder := json.NewDecoder(r.file)
	u := &model.URL{}
	for {
		if err := decoder.Decode(u); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error decoding file : %v", err)
		}
		url := &model.URL{ShortURL: u.ShortURL, LongURL: u.LongURL}
		r.cacheByShort[u.ShortURL] = url
		r.cacheByUser[u.UserID] = append(r.cacheByUser[u.UserID], url)
	}
	return nil
}

// Get returns original link by id or an error if id is not present
func (r *FileRepo) Get(short string, ctx context.Context) (string, error) {
	r.Lock()
	defer r.Unlock()
	url, ok := r.cacheByShort[short]
	if !ok {
		return "", fmt.Errorf("invalid id: %s", short)
	}
	return url.LongURL, nil
}

// Add adds a link to db and returns assigned id
func (r *FileRepo) Add(url *model.URL, ctx context.Context) (bool, error) {
	r.Lock()
	defer r.Unlock()
	if v, k := r.existingURLs[url.LongURL]; k {
		url.ShortURL = v.ShortURL
		return true, nil
	}
	r.cacheByShort[url.ShortURL] = url
	r.cacheByUser[url.UserID] = append(r.cacheByUser[url.UserID], url)
	r.existingURLs[url.LongURL] = url
	err := r.encoder.Encode(*url)
	if err != nil {
		return false, err
	}
	return false, nil
}

func (r *FileRepo) AddBatch(urls map[string]*model.URL, ctx context.Context) (bool, error) {
	r.Lock()
	defer r.Unlock()
	var duplicates bool
	for _, v := range urls {
		if i, k := r.existingURLs[v.LongURL]; k {
			duplicates = true
			v.ShortURL = i.ShortURL
			continue
		}
		r.cacheByShort[v.ShortURL] = v
		r.cacheByUser[v.UserID] = append(r.cacheByUser[v.UserID], v)
		err := r.encoder.Encode(&v)
		if err != nil {
			return false, err
		}
	}
	return duplicates, nil
}

// NewID returns a number to encode as an id
func (r *FileRepo) NewID() (int, error) {
	r.Lock()
	defer r.Unlock()
	r.lastID++
	return r.lastID, nil
}

func (r *FileRepo) GetByUser(userID string, ctx context.Context) ([]*model.URL, error) {
	r.RLock()
	defer r.RUnlock()
	var s []*model.URL
	s = append(s, r.cacheByUser[userID]...)
	if len(s) == 0 {
		return nil, nil
	}
	return s, nil
}

func (r *FileRepo) Ping(ctx context.Context) error {
	_, err := os.Stat(r.file.Name())
	return err
}

func (r *FileRepo) DeleteRepo(ctx context.Context) error {
	err := r.file.Close()
	if err != nil {
		return fmt.Errorf("error closing file : %v", err)
	}
	err = os.Remove(r.file.Name())
	if err != nil {
		return fmt.Errorf("error deleting file : %v", err)
	}
	return nil
}
