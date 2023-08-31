package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/router/loader"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/validators"
)

// ShortenerImpl is a ShortenerService implementation
type ShortenerImpl struct {
	repo   storage.Repository
	cfg    *config.Config
	loader *loader.UserLoader
}

// NewShortenerImpl returns a ShortenerImpl implementation
func NewShortenerImpl(repo storage.Repository, cfg *config.Config) *ShortenerImpl {
	return &ShortenerImpl{
		repo:   repo,
		cfg:    cfg,
		loader: loader.NewDeleteLoader(repo),
	}
}

// Expand gets original url from short
func (s *ShortenerImpl) Expand(ctx context.Context, id string) (*models.URL, error) {
	// If the URL has been deleted, return Gone status.
	url, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", models.ErrRepoError, err.Error())
	}
	if url.Deleted {
		return url, models.ErrURLDeleted
	}
	return url, nil
}

// APIDeleteBatch deletes a batch of urls by user
func (s *ShortenerImpl) DeleteBatch(urlIDs []string, userID string) {
	// Create a slice of DeleteURLItem objects from the URL IDs.
	deleteURLs := make([]*models.DeleteURLItem, len(urlIDs))
	for i, v := range urlIDs {
		deleteURLs[i] = &models.DeleteURLItem{UserID: userID, ShortURL: v}
	}
	// Start a goroutine to delete the URLs asynchronously.
	go func() {
		num, err := s.loader.LoadAll(deleteURLs)
		if err[0] != nil {
			log.Printf("error deleing urls :%v", err[0])
		}
		var result int
		for _, v := range num {
			result += v
		}
		log.Printf("deleted %v urls", result)
	}()
}

// Ping checks availibility
func (s *ShortenerImpl) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

// Shorten shortens a url
func (s *ShortenerImpl) Shorten(ctx context.Context, url *models.URL) (*models.URL, error) {
	// Check if body is a valid URL.
	if !validators.IsURL(url.LongURL) {
		return nil, models.ErrInvalidURL
	}
	var err error
	// Get short url for long
	id, err := s.repo.NewID(url.LongURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", models.ErrRepoError, err.Error())
	}
	url.ShortURL = id
	// Add record to repo
	duplicates, err := s.repo.Add(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", models.ErrRepoError, err.Error())
	}
	if duplicates {
		return url, models.ErrDuplicate
	}
	return url, nil
}

// ShortenBatch shortens multiple urls
func (s *ShortenerImpl) ShortenBatch(ctx context.Context, userID string, urls []*models.URL) ([]*models.URL, error) {
	var err error
	for i, v := range urls {
		// Check if the original URL is valid.
		if !validators.IsURL(v.LongURL) {
			// If the URL is not valid, set the response item to indicate a bad URL request.
			urls[i].ShortURL = "incorrect url"
			continue
		}
		// Generate a short ID for the URL.
		urls[i].ShortURL, err = s.repo.NewID(v.LongURL)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", models.ErrRepoError, err.Error())
		}
		// Add short url to info.
	}
	// Add the URLs to the repository.
	duplicates, err := s.repo.AddBatch(ctx, urls)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", models.ErrRepoError, err.Error())
	}
	if duplicates {
		return urls, models.ErrDuplicate
	}
	return urls, nil
}

// Stats gets the count of urls and registered users
func (s *ShortenerImpl) Stats(ctx context.Context) (*models.Stats, error) {
	return s.repo.Stats(ctx)
}

// ExpandUser gets user links
func (s *ShortenerImpl) ExpandUser(ctx context.Context, userID string) ([]*models.URL, error) {
	// Get the list of URLs created by user.
	urls, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(urls) == 0 {
		return nil, models.ErrNoContent
	}
	return urls, nil
}

// BuildURL appends domain to a short link when using rest api
func (s *ShortenerImpl) BuildURL(url string) string {
	return strings.Join([]string{s.cfg.BaseURL, url}, "/")
}
