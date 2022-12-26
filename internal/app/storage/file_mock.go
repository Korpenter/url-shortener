package storage

import (
	"encoding/json"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"io"
	"log"
	"os"

	"github.com/Mldlr/url-shortener/internal/app/models"
)

type mockFileRepo struct {
	file         *os.File
	cacheByShort map[string]*models.URL
	cacheByUser  map[string][]*models.URL
	encoder      json.Encoder
}

// NewMockFileRepo itiates new mock file repo, creating a file and adding a record to it
func newMockFileRepo() (*mockFileRepo, error) {
	urls := []models.URL{
		{
			ShortURL: "1",
			LongURL:  "yandex.ru",
			UserID:   "Helloworld",
		},
		{
			ShortURL: "2",
			LongURL:  "hero.ru",
			UserID:   "Helloworld",
		},
	}
	file, err := os.OpenFile("./mockFileDB", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("error creating mock file : %v", err)
	}
	mock := mockFileRepo{
		file:         file,
		cacheByShort: make(map[string]*models.URL),
		cacheByUser:  make(map[string][]*models.URL),
		encoder:      *json.NewEncoder(file),
	}
	for _, v := range urls {
		if _, err = mock.add(v.LongURL, v.ShortURL, v.UserID); err != nil {
			return nil, fmt.Errorf("error adding mock records : %v", err)
		}
	}
	return &mock, nil
}

// DeleteMock deletes mock file
func (r *mockFileRepo) delete() error {
	err := r.file.Close()
	if err != nil {
		return fmt.Errorf("error closing mock file : %v", err)
	}
	err = os.Remove("./mockFileDB")
	if err != nil {
		return fmt.Errorf("error deleting mock file : %v", err)
	}
	return nil
}

// Load loads stored url records from file
func (r *mockFileRepo) load() error {
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
func (r *mockFileRepo) get(short string) (*models.URL, error) {
	url, ok := r.cacheByShort[short]
	if !ok {
		return nil, fmt.Errorf("invalid id: %s", short)
	}
	return url, nil
}

// Add adds a link to db and returns assigned id
func (r *mockFileRepo) add(longURL, short, userID string) (string, error) {
	url := &models.URL{ShortURL: short, LongURL: longURL, UserID: userID}
	r.cacheByShort[short] = url
	r.cacheByUser[userID] = append(r.cacheByUser[userID], url)
	err := r.encoder.Encode(*url)
	if err != nil {
		return short, err
	}
	return short, nil
}

func (r *mockFileRepo) addBatch(urls []models.URL) error {
	for _, v := range urls {
		r.cacheByShort[v.ShortURL] = &v
		r.cacheByUser[v.UserID] = append(r.cacheByUser[v.UserID], &v)
		err := r.encoder.Encode(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *mockFileRepo) updateFile() {
	log.Println("starting file update")
	err := r.file.Truncate(0)
	if err != nil {
		log.Println("error truncating file:", err)
	}
	_, err = r.file.Seek(0, 0)
	if err != nil {
		log.Println("error setting pointer in file:", err)
	}
	for _, v := range r.cacheByShort {
		err = r.encoder.Encode(&v)
		if err != nil {
			log.Println("error encoding url:", err)
		}
	}
	log.Println("finished file update")
}

func (r *FileRepo) newID(url string) (string, error) {
	return encoders.ToRBase62(url), nil
}

func (r *mockFileRepo) getByUser(userID string) ([]*models.URL, error) {
	urls := make([]*models.URL, 0)
	urls = append(urls, r.cacheByUser[userID]...)
	if len(urls) == 0 {
		return nil, fmt.Errorf("no urls found for user")
	}
	return urls, nil
}

func (r *mockFileRepo) ping() error {
	_, err := os.Stat(r.file.Name())
	return err
}
