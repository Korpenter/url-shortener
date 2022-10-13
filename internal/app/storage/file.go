package storage

import (
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
	encoder      json.Encoder
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
		encoder:      *json.NewEncoder(file),
	}, nil
}

// Load loads stored url records from file
func (r *FileRepo) Load() error {
	decoder := json.NewDecoder(r.file)
	u := &model.URL{}
	for {
		if err := decoder.Decode(u); err == io.EOF {
			fmt.Println(err)
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
func (r *FileRepo) Get(short string) (string, error) {
	r.Lock()
	defer r.Unlock()
	url, ok := r.cacheByShort[short]
	if !ok {
		return "", fmt.Errorf("invalid id: %s", short)
	}
	return url.LongURL, nil
}

// Add adds a link to db and returns assigned id
func (r *FileRepo) Add(url *model.URL) error {
	r.Lock()
	defer r.Unlock()
	r.cacheByShort[url.ShortURL] = url
	r.cacheByUser[url.UserID] = append(r.cacheByUser[url.UserID], url)
	err := r.encoder.Encode(*url)
	if err != nil {
		return err
	}
	return nil
}

func (r *FileRepo) AddBatch(urls []*model.URL) error {
	r.Lock()
	defer r.Unlock()
	for _, v := range urls {
		r.cacheByShort[v.ShortURL] = v
		r.cacheByUser[v.UserID] = append(r.cacheByUser[v.UserID], v)
		err := r.encoder.Encode(&v)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewID returns a number to encode as an id
func (r *FileRepo) NewID() (int, error) {
	r.Lock()
	defer r.Unlock()
	return len(r.cacheByShort) + 1, nil
}

func (r *FileRepo) GetByUser(userID string) ([]*model.URL, error) {
	r.RLock()
	defer r.RUnlock()
	s := []*model.URL{}
	s = append(s, r.cacheByUser[userID]...)
	if len(s) == 0 {
		return nil, nil
	}
	return s, nil
}

func (r *FileRepo) Ping() error {
	_, err := os.Stat(r.file.Name())
	return err
}
