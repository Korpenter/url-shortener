package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// MockFileRepo is a mock of a FileRepo
type MockFileRepo struct {
	file    *os.File
	cache   map[string]string
	encoder json.Encoder
}

// NewMockFileRepo itiates new mock file repo, creating a file and adding a record to it
func NewMockFileRepo() (*MockFileRepo, error) {
	urls := []url{
		{
			ID:      "1",
			LongURL: "yandex.ru",
		},
		{
			ID:      "2",
			LongURL: "hero.ru",
		},
	}
	file, err := os.OpenFile("./mockFileDB", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("error creating mock file : %v", err)
	}
	mock := MockFileRepo{
		file:    file,
		cache:   make(map[string]string),
		encoder: *json.NewEncoder(file),
	}
	for _, v := range urls {
		if _, err = mock.Add(v.LongURL, v.ID); err != nil {
			return nil, fmt.Errorf("error adding mock records : %v", err)
		}
	}
	return &mock, nil
}

// DeleteMock deletes mock file
func (r *MockFileRepo) DeleteMock() error {
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
func (r *MockFileRepo) Load() error {
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
func (r *MockFileRepo) Get(id string) (string, error) {
	longURL, ok := r.cache[id]
	if !ok {
		return "", fmt.Errorf("invalid id: %s", id)
	}
	return longURL, nil
}

// Add adds a link to db and returns assigned id
func (r *MockFileRepo) Add(longURL, id string) (string, error) {
	r.cache[id] = longURL
	url := url{
		ID:      id,
		LongURL: longURL,
	}
	err := r.encoder.Encode(url)
	if err != nil {
		return id, err
	}
	return id, nil
}

// NewID returns a number to encode as an id
func (r *MockFileRepo) NewID() (int, error) {
	return len(r.cache) + 1, nil
}
