package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/Mldlr/url-shortener/internal/app/utils/helpers"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepo is a Postgres db storage.
type PostgresRepo struct {
	conn *pgxpool.Pool
	sync.Mutex
}

// NewPostgresRepo initializes Postgres DB storage instance from connection string.
func NewPostgresRepo(connString string) (*PostgresRepo, error) {
	// Create config from connection string.
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	// Create connection pool from config.
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	return &PostgresRepo{conn: conn}, nil
}

// NewTableURLs creates the 'urls' table in the database if it does not already exist.
func (r *PostgresRepo) NewTableURLs() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := r.conn.Exec(ctx, createUrls)
	if err != nil {
		return err
	}
	return nil
}

// Get returns original link by id or an error if id is not present
func (r *PostgresRepo) Get(ctx context.Context, id string) (*models.URL, error) {
	var url models.URL
	err := r.conn.QueryRow(ctx, getQuery, id).Scan(&url.LongURL, &url.Deleted)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", id)
	}
	return &url, nil
}

// GetByUser finds URLs created by a specific user.
func (r *PostgresRepo) GetByUser(ctx context.Context, userID string) ([]*models.URL, error) {
	var url models.URL
	var count int
	err := r.conn.QueryRow(ctx, countUserURLs, userID).Scan(&count)
	if err != nil {
		return nil, err
	}
	// Preallocate a slice of URL
	urls := make([]*models.URL, 0, count)
	rows, err := r.conn.Query(ctx, getByUserQuery, userID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()
	// For each row return read values to structure and append to URL slice
	for rows.Next() {
		err = rows.Scan(&url.ShortURL, &url.LongURL, &url.UserID, &url.Deleted)
		if err != nil {
			return nil, err
		}
		urls = append(urls, &url)
	}
	return urls, nil
}

// Add adds a link to db and returns assigned id
func (r *PostgresRepo) Add(ctx context.Context, url *models.URL) (bool, error) {
	var duplicates bool
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer helpers.CommitTx(ctx, tx, err)
	// Execute insert query and read inserted ID.
	err = tx.QueryRow(ctx, addQuery, url.ShortURL, url.LongURL, url.UserID).Scan(&url.ShortURL)
	if err != nil {
		// If no rows inserted, but query was successful.
		if errors.Is(err, pgx.ErrNoRows) {
			duplicates = true
			// Select existing id.
			err = tx.QueryRow(ctx, getShort, url.LongURL).Scan(&url.ShortURL)
		} else {
			return false, err
		}
	}
	return duplicates, err
}

// AddBatch adds multiple URLs to repository.
func (r *PostgresRepo) AddBatch(ctx context.Context, urls map[string]*models.URL) (bool, error) {
	var duplicates bool
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer helpers.CommitTx(ctx, tx, err)
	// For every URL.
	for _, v := range urls {
		// Execute insert query and read inserted ID.
		err = tx.QueryRow(ctx, addQuery, v.ShortURL, v.LongURL, v.UserID).Scan(&v.ShortURL)
		if err != nil {
			// If no rows inserted, but query was successful.
			if errors.Is(err, pgx.ErrNoRows) {
				duplicates = true
				// Select existing id.
				err = tx.QueryRow(ctx, getShort, v.LongURL).Scan(&v.ShortURL)
			}
		}
	}
	return duplicates, err
}

// DeleteURLs delete urls from cache.
func (r *PostgresRepo) DeleteURLs(deleteURLs []*models.DeleteURLItem) (int, error) {
	ctx := context.Background()
	tx, err := r.conn.Begin(ctx)
	var n int
	if err != nil {
		return n, err
	}
	defer helpers.CommitTx(ctx, tx, err)
	// Preallocate slice to store short URLs.
	shortURLs := make([]string, len(deleteURLs))
	for i, v := range deleteURLs {
		shortURLs[i] = v.ShortURL
	}
	// Update URLs.
	res, err := tx.Exec(ctx, updateDeleteQuery, shortURLs, deleteURLs[0].UserID)
	if err != nil {
		return n, err
	}
	n = int(res.RowsAffected())
	return n, err
}

// NewID calculates a string to use as an ID.
func (r *PostgresRepo) NewID(url string) (string, error) {
	return encoders.ToRBase62(url), nil
}

// Ping checks if file is available.
func (r *PostgresRepo) Ping(ctx context.Context) error {
	return r.conn.Ping(ctx)
}

// DeleteRepo deletes repository tables.
func (r *PostgresRepo) DeleteRepo(ctx context.Context) error {
	// drop tables
	_, err := r.conn.Exec(ctx, drop)
	if err != nil {
		return err
	}
	return nil
}
