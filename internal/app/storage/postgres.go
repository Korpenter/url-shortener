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
	"sync"
)

type PostgresRepo struct {
	conn *pgxpool.Pool
	sync.Mutex
}

func NewPostgresRepo(connString string) (*PostgresRepo, error) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	return &PostgresRepo{conn: conn}, nil
}

func (r *PostgresRepo) NewTableURLs() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := r.conn.Exec(ctx, createUrls)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepo) Get(ctx context.Context, id string) (*models.URL, error) {
	var url models.URL
	err := r.conn.QueryRow(ctx, getQuery, id).Scan(&url.LongURL, &url.Deleted)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", id)
	}
	return &url, nil
}

func (r *PostgresRepo) GetByUser(ctx context.Context, userID string) ([]*models.URL, error) {
	var url models.URL
	var count int
	err := r.conn.QueryRow(ctx, countUserURLs, userID).Scan(&count)
	if err != nil {
		return nil, err
	}
	urls := make([]*models.URL, 0, count)
	rows, err := r.conn.Query(ctx, getByUserQuery, userID)
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

func (r *PostgresRepo) Add(ctx context.Context, url *models.URL) (bool, error) {
	var duplicates bool
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer helpers.CommitTx(ctx, tx, err)
	err = tx.QueryRow(ctx, addQuery, url.ShortURL, url.LongURL, url.UserID).Scan(&url.ShortURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			duplicates = true
			err = tx.QueryRow(ctx, getShort, url.LongURL).Scan(&url.ShortURL)
		} else {
			return false, err
		}
	}
	return duplicates, err
}

func (r *PostgresRepo) AddBatch(ctx context.Context, urls map[string]*models.URL) (bool, error) {
	var duplicates bool
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer helpers.CommitTx(ctx, tx, err)
	for _, v := range urls {
		err = tx.QueryRow(ctx, addQuery, v.ShortURL, v.LongURL, v.UserID).Scan(&v.ShortURL)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				duplicates = true
				err = tx.QueryRow(ctx, getShort, v.LongURL).Scan(&v.ShortURL)
			}
		}
	}
	return duplicates, err
}

func (r *PostgresRepo) DeleteURLs(deleteURLs []*models.DeleteURLItem) (int, error) {
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
	res, err := tx.Exec(ctx, updateDeleteQuery, shortURLs, deleteURLs[0].UserID)
	if err != nil {
		return n, err
	}
	n = int(res.RowsAffected())
	return n, err
}

func (r *PostgresRepo) NewID(url string) (string, error) {
	return encoders.ToRBase62(url), nil
}

func (r *PostgresRepo) Ping(ctx context.Context) error {
	return r.conn.Ping(ctx)
}

func (r *PostgresRepo) DeleteRepo(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, drop)
	if err != nil {
		return err
	}
	return nil
}
