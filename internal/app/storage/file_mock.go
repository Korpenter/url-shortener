package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Mldlr/url-shortener/internal/app/model"
)

// MockFileRepo is a mock of a FileRepo
type mockFileRepo struct {
	file         *os.File
	cacheByShort map[string]*model.URL
	cacheByUser  map[string][]*model.URL
	encoder      json.Encoder
}

// NewMockFileRepo itiates new mock file repo, creating a file and adding a record to it
func newMockFileRepo() (*mockFileRepo, error) {
	urls := []model.URL{
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
		cacheByShort: make(map[string]*model.URL),
		cacheByUser:  make(map[string][]*model.URL),
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
func (r *mockFileRepo) deleteMock() error {
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
func (r *mockFileRepo) get(short string) (string, error) {
	url, ok := r.cacheByShort[short]
	if !ok {
		return "", fmt.Errorf("invalid id: %s", short)
	}
	return url.LongURL, nil
}

// Add adds a link to db and returns assigned id
func (r *mockFileRepo) add(longURL, short, userID string) (string, error) {
	url := &model.URL{ShortURL: short, LongURL: longURL, UserID: userID}
	r.cacheByShort[short] = url
	r.cacheByUser[userID] = append(r.cacheByUser[userID], url)
	err := r.encoder.Encode(*url)
	if err != nil {
		return short, err
	}
	return short, nil
}

func (r *mockFileRepo) getByUser(userID string) ([]*model.URL, error) {
	s := []*model.URL{}
	s = append(s, r.cacheByUser[userID]...)
	if len(s) == 0 {
		return nil, fmt.Errorf("no urls found for user")
	}
	return s, nil
}

func (r *mockFileRepo) ping() error {
	_, err := os.Stat(r.file.Name())
	return err
}
