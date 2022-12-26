package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"io"
	"log"
	"os"
	"sync"
)

// FileRepo is an in-file url storage
type FileRepo struct {
	file         *os.File
	cacheByShort map[string]*models.URL
	cacheByUser  map[string][]*models.URL
	existingURLs map[string]*models.URL
	encoder      json.Encoder
	sync.RWMutex
}

// NewFileRepo returns a pointer to a new repo instance
func NewFileRepo(filename string) (*FileRepo, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("error openin file : %v", err)
	}
	return &FileRepo{
		file:         file,
		cacheByShort: make(map[string]*models.URL),
		cacheByUser:  make(map[string][]*models.URL),
		existingURLs: make(map[string]*models.URL),
		encoder:      *json.NewEncoder(file),
	}, nil
}

// Load loads stored url records from file
func (r *FileRepo) Load() error {
	decoder := json.NewDecoder(r.file)
	u := &models.URL{}
	for {
		if err := decoder.Decode(u); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error decoding file : %v", err)
		}
		url := &models.URL{ShortURL: u.ShortURL, LongURL: u.LongURL}
		r.cacheByShort[u.ShortURL] = url
		r.cacheByUser[u.UserID] = append(r.cacheByUser[u.UserID], url)
	}
	return nil
}

// Get returns original link by id or an error if id is not present
func (r *FileRepo) Get(ctx context.Context, id string) (*models.URL, error) {
	r.Lock()
	defer r.Unlock()
	url, ok := r.cacheByShort[id]
	if !ok {
		return nil, fmt.Errorf("invalid id: %s", id)
	}
	return url, nil
}

// Add adds a link to db and returns assigned id
func (r *FileRepo) Add(ctx context.Context, url *models.URL) (bool, error) {
	r.Lock()
	defer r.Unlock()
	if v, k := r.existingURLs[url.LongURL]; k {
		url.ShortURL = v.ShortURL
		return true, nil
	}
	r.cacheByShort[url.ShortURL] = url
	r.cacheByUser[url.UserID] = append(r.cacheByUser[url.UserID], url)
	r.existingURLs[url.LongURL] = url
	return false, nil
}

func (r *FileRepo) AddBatch(ctx context.Context, urls map[string]*models.URL) (bool, error) {
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
	}
	return duplicates, nil
}

// NewID returns a number to encode as an id
func (r *FileRepo) NewID(url string) (string, error) {
	return encoders.ToRBase62(url), nil
}

func (r *FileRepo) GetByUser(ctx context.Context, userID string) ([]*models.URL, error) {
	r.RLock()
	defer r.RUnlock()
	var s []*models.URL
	s = append(s, r.cacheByUser[userID]...)
	if len(s) == 0 {
		return nil, nil
	}
	return s, nil
}

func (r *FileRepo) DeleteURLs(deleteURLs []*models.DeleteURLItem) (int, error) {
	r.Lock()
	defer r.Unlock()
	var n int
	for _, v := range deleteURLs {
		if r.cacheByShort[v.ShortURL].UserID == v.UserID {
			r.cacheByShort[v.ShortURL].Deleted = true
			n++
		}
	}
	return n, nil
}

func (r *FileRepo) updateFile() {
	r.Lock()
	defer r.Unlock()
	log.Println("starting file update")
	err := r.file.Truncate(0)
	if err != nil {
		log.Println("error updating file 1:", err)
	}
	_, err = r.file.Seek(0, 0)
	if err != nil {
		log.Println("error updating file 2:", err)
	}
	for _, v := range r.cacheByShort {
		err = r.encoder.Encode(&v)
		if err != nil {
			log.Println("error updating file 3:", err)
		}
	}
	log.Println("finished file update")
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
