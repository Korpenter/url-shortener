package models

import "errors"

// Declaration of errors for shortener
var (
	// ErrURLNotFound - url not registered
	ErrURLNotFound = errors.New("URL not found")
	// ErrURLDeleted - ulr deleted
	ErrURLDeleted = errors.New("URL deleted")
	// ErrInvalidURL - invalid url
	ErrInvalidURL = errors.New("invalid url")
	// ErrGettingID - repo error
	ErrRepoError = errors.New("repository error")
	// ErrDuplicate - duplicate long url
	ErrDuplicate = errors.New("url already in db")
	// ErrNoContent - no registered url for user
	ErrNoContent = errors.New("no urls for user")
)
