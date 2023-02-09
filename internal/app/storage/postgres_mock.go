package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/Mldlr/url-shortener/internal/app/utils/helpers"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	mockCreateUrls = `CREATE TABLE IF NOT EXISTS urls_test (
            	short varchar(255) PRIMARY KEY,
                original varchar(255),
    			userid varchar(64),
    			deleted boolean DEFAULT false,
    			UNIQUE(original)
                )`
	mockAddQuery = `
	INSERT INTO urls_test (short, original, userid)
	VALUES ($1, $2, $3)
	ON CONFLICT DO NOTHING
	RETURNING short`
	mockUpdateDeleteQuery = `UPDATE urls_test SET deleted=TRUE WHERE short IN (SELECT unnest($1::text[])) AND userid = $2`
	mockGetQuery          = `SELECT original, deleted FROM urls_test WHERE short = $1`
	mockGetByUserQuery    = `SELECT * FROM urls_test WHERE userid = $1`
	mockGetShort          = `SELECT short FROM urls_test WHERE original = $1`
	mockCountUserURLs     = "SELECT count(*) FROM urls_test WHERE userid = $1"
	getMockStats          = "SELECT COUNT(*), COUNT(DISTINCT(userid)) FROM urls_test;"
	mockDrop              = `DROP TABLE urls_test`
)

type postgresMockRepo struct {
	conn *pgxpool.Pool
}

// NewPostgresMockRepo initializes Postgres DB test storage instance from connection string with data.
func NewPostgresMockRepo(connString string) (*postgresMockRepo, error) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	url1 := &models.URL{ShortURL: "1", LongURL: "https://github.com/Mldlr/url-shortener/internal/app/utils/encoders"}
	url2 := &models.URL{ShortURL: "2", LongURL: "https://yandex.ru/"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = conn.Exec(ctx, mockCreateUrls)
	if err != nil {
		return nil, err
	}
	repo := &postgresMockRepo{conn: conn}
	repo.Add(context.Background(), url1)
	repo.Add(context.Background(), url2)
	return repo, nil
}

// Get returns original link by id or an error if id is not present.
func (r *postgresMockRepo) Get(ctx context.Context, id string) (*models.URL, error) {
	var url models.URL
	err := r.conn.QueryRow(ctx, mockGetQuery, id).Scan(&url.LongURL, &url.Deleted)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", id)
	}
	return &url, nil
}

// GetByUser finds URLs created by a specific user.
func (r *postgresMockRepo) GetByUser(ctx context.Context, userID string) ([]*models.URL, error) {
	var url models.URL
	var count int
	err := r.conn.QueryRow(ctx, mockCountUserURLs, userID).Scan(count)
	if err != nil {
		return nil, err
	}
	urls := make([]*models.URL, 0, count)
	rows, err := r.conn.Query(ctx, mockGetByUserQuery, userID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()
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
func (r *postgresMockRepo) Add(ctx context.Context, url *models.URL) (bool, error) {
	var duplicates bool
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer helpers.CommitTx(ctx, tx, err)
	err = tx.QueryRow(ctx, mockAddQuery, url.ShortURL, url.LongURL, url.UserID).Scan(&url.ShortURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			duplicates = true
			err = tx.QueryRow(ctx, mockGetShort, url.LongURL).Scan(&url.ShortURL)
		}
	}
	return duplicates, err
}

// AddBatch adds multiple URLs to repository.
func (r *postgresMockRepo) AddBatch(ctx context.Context, urls map[string]*models.URL) (bool, error) {
	var duplicates bool
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer helpers.CommitTx(ctx, tx, err)
	for _, v := range urls {
		err = tx.QueryRow(ctx, mockAddQuery, v.ShortURL, v.LongURL, v.UserID).Scan(&v.ShortURL)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				duplicates = true
				err = tx.QueryRow(ctx, mockGetShort, v.LongURL).Scan(&v.ShortURL)
			}
		}
	}
	return duplicates, err
}

// DeleteURLs delete urls from cache.
func (r *postgresMockRepo) NewID(url string) (string, error) {
	return encoders.ToRBase62(url), nil
}

// NewID calculates a string to use as an ID.
func (r *postgresMockRepo) DeleteURLs(deleteURLs []*models.DeleteURLItem) (int, error) {
	ctx := context.Background()
	tx, err := r.conn.Begin(ctx)
	var n int
	if err != nil {
		return n, err
	}
	defer helpers.CommitTx(ctx, tx, err)
	shortURLs := make([]string, len(deleteURLs))
	for i, v := range deleteURLs {
		shortURLs[i] = v.ShortURL
	}
	res, err := tx.Exec(ctx, mockUpdateDeleteQuery, shortURLs, deleteURLs[0].UserID)
	if err != nil {
		return n, err
	}
	n = int(res.RowsAffected())
	return n, err
}

// Ping checks if file is available.
func (r *postgresMockRepo) Ping(ctx context.Context) error {
	return r.conn.Ping(ctx)
}

// Stats gets count of urls and registered users
func (r *postgresMockRepo) Stats(ctx context.Context) (*models.Stats, error) {
	var stats models.Stats
	err := r.conn.QueryRow(ctx, getMockStats).Scan(&stats.URLCount, &stats.UserCount)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// DeleteRepo deletes repository tables.
func (r *postgresMockRepo) DeleteRepo(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, mockDrop)
	if err != nil {
		return err
	}
	return nil
}

// Close closes connection to db
func (r *postgresMockRepo) Close() error {
	r.conn.Close()
	return nil
}
