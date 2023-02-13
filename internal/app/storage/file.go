package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
)

// FileRepo is an in-file url storage
type FileRepo struct {
	// file stores URL data.
	file *os.File
	// cacheByShort maps short URLs to their corresponding URL models.
	cacheByShort map[string]*models.URL
	// cacheByUser maps user IDs to lists of URL models.
	cacheByUser map[string][]*models.URL
	// existingURLs maps long URLs to their corresponding URL models.
	existingURLs map[string]*models.URL
	// encoder encodes URL data for storage in the file.
	encoder json.Encoder
	// RWMutex synchronizes access to the FileRepo.
	sync.RWMutex
}

// NewFileRepo initializes a new in-file storage.
func NewFileRepo(filename string) (*FileRepo, error) {
	// Create or open file
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

// Load loads stored url records from file.
func (r *FileRepo) Load() error {
	// Decode file
	decoder := json.NewDecoder(r.file)
	u := &models.URL{}
	for {
		if err := decoder.Decode(u); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error decoding file : %v", err)
		}
		// Add decoded URL to maps
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
	// check for URL in map
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
	// check for url in map and return if it already exists
	if v, k := r.existingURLs[url.LongURL]; k {
		url.ShortURL = v.ShortURL
		return true, nil
	}
	// otherwise add url to maps
	r.cacheByShort[url.ShortURL] = url
	r.cacheByUser[url.UserID] = append(r.cacheByUser[url.UserID], url)
	r.existingURLs[url.LongURL] = url
	return false, nil
}

// AddBatch adds multiple URLs to repository.
func (r *FileRepo) AddBatch(ctx context.Context, urls []*models.URL) (bool, error) {
	r.Lock()
	defer r.Unlock()
	var duplicates bool
	// for each url check if url is in map and add it otherwise
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

// NewID calculates a string to use as an ID.
func (r *FileRepo) NewID(url string) (string, error) {
	return encoders.ToRBase62(url), nil
}

// GetByUser finds URLs created by a specific user.
func (r *FileRepo) GetByUser(ctx context.Context, userID string) ([]*models.URL, error) {
	r.RLock()
	defer r.RUnlock()
	s := make([]*models.URL, 0)
	// Get all urls for user
	s = append(s, r.cacheByUser[userID]...)
	if len(s) == 0 {
		return nil, nil
	}
	return s, nil
}

// DeleteURLs delete urls from cache.
func (r *FileRepo) DeleteURLs(deleteURLs []*models.DeleteURLItem) (int, error) {
	r.Lock()
	defer r.Unlock()
	var n int
	// For each of the urls check if the user created this url and delete it if confirmed
	for _, v := range deleteURLs {
		if r.cacheByShort[v.ShortURL].UserID == v.UserID {
			r.cacheByShort[v.ShortURL].Deleted = true
			n++
		}
	}
	return n, nil
}

func (r *FileRepo) update() {
	r.Lock()
	defer r.Unlock()
	r.updateFile()
}

// updateFile updates file contents from cache.
func (r *FileRepo) updateFile() {
	log.Println("starting file update")
	// Truncate file
	err := r.file.Truncate(0)
	if err != nil {
		log.Println("error updating file 1:", err)
	}
	// Put cursor at the beginning of the file
	_, err = r.file.Seek(0, 0)
	if err != nil {
		log.Println("error updating file 2:", err)
	}
	// Write cached data to file
	for _, v := range r.cacheByShort {
		err = r.encoder.Encode(&v)
		if err != nil {
			log.Println("error updating file 3:", err)
		}
	}
	log.Println("finished file update")
}

// Ping checks if file is available.
func (r *FileRepo) Ping(ctx context.Context) error {
	_, err := os.Stat(r.file.Name())
	return err
}

// DeleteRepo deletes repository file.
func (r *FileRepo) DeleteRepo(ctx context.Context) error {
	err := r.file.Close()
	if err != nil {
		return fmt.Errorf("error closing file : %v", err)
	}
	// Delete the file
	err = os.Remove(r.file.Name())
	if err != nil {
		return fmt.Errorf("error deleting file : %v", err)
	}
	return nil
}

// Stats gets count of urls and registered users
func (r *FileRepo) Stats(ctx context.Context) (*models.Stats, error) {
	r.RLock()
	defer r.RUnlock()
	var stats models.Stats
	stats.URLCount = len(r.existingURLs)
	stats.UserCount = len(r.cacheByUser)
	return &stats, nil
}

// Close closes file
func (r *FileRepo) Close() error {
	r.Lock()
	defer r.Unlock()
	r.updateFile()
	err := r.file.Close()
	if err != nil {
		return fmt.Errorf("error closing file : %v", err)
	}
	return nil
}
