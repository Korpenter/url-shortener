package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

// FileRepo is an in-file url storage
type FileRepo struct {
	file  *os.File
	cache map[string]string
	sync.Mutex
}

// NewFileRepo returns a pointer to a new repo instance
func NewFileRepo(filename string) (*FileRepo, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error openin file : %v", err)
	}
	return &FileRepo{
		file:  file,
		cache: make(map[string]string),
	}, nil
}

// Load loads stored url records from file
func (r *FileRepo) Load() error {
	decoder := json.NewDecoder(r.file)
	u := &url{}
	for {
		if err := decoder.Decode(u); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error decoding file : %v", err)
		}
		r.cache[u.ID] = u.LongURL
	}
	return nil
}

// Get returns original link by id or an error if id is not present
func (r *FileRepo) Get(id string) (string, error) {
	r.Lock()
	defer r.Unlock()
	longURL, ok := r.cache[id]
	if !ok {
		return "", fmt.Errorf("invalid id: %s", id)
	}
	return longURL, nil
}

// Add adds a link to db and returns assigned id
func (r *FileRepo) Add(longURL, id string) (string, error) {
	r.Lock()
	defer r.Unlock()
	r.cache[id] = longURL
	url := url{
		ID:      id,
		LongURL: longURL,
	}
	encoder := json.NewEncoder(r.file)
	err := encoder.Encode(url)
	if err != nil {
		return id, err
	}
	return id, nil
}

// NewID returns a number to encode as an id
func (r *FileRepo) NewID() (int, error) {
	r.Lock()
	defer r.Unlock()
	return len(r.cache) + 1, nil
}
