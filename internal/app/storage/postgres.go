package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

type PostgresRepo struct {
	conn   *pgxpool.Pool
	lastID int
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

func (r *PostgresRepo) NumberOfURLs() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rows, err := r.conn.Query(ctx, countURL)
	if err != nil {
		return err
	}
	for rows.Next() {
		err = rows.Scan(&r.lastID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresRepo) Get(id string, ctx context.Context) (string, error) {
	var url string
	err := r.conn.QueryRow(ctx, getQuery, id).Scan(&url)
	if err != nil {
		return "", fmt.Errorf("invalid id: %v", id)
	}
	return url, nil
}

func (r *PostgresRepo) GetByUser(userID string, ctx context.Context) ([]*model.URL, error) {
	var url model.URL
	urls := make([]*model.URL, 0)
	rows, err := r.conn.Query(ctx, getByUserQuery, userID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&url.ShortURL, &url.LongURL, &url.UserID)
		if err != nil {
			return nil, err
		}
		urls = append(urls, &url)
	}
	return urls, nil
}

func (r *PostgresRepo) Add(url *model.URL, ctx context.Context) (bool, error) {
	var duplicates bool
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	err = tx.QueryRow(ctx, addQuery, url.ShortURL, url.LongURL, url.UserID).Scan(&url.ShortURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			duplicates = true
			err = tx.QueryRow(ctx, getShort, url.LongURL).Scan(&url.ShortURL)
		}
		if err != nil {
			return false, err
		}
	}
	return duplicates, nil
}

func (r *PostgresRepo) AddBatch(urls map[string]*model.URL, ctx context.Context) (bool, error) {
	var duplicates bool
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
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

func (r *PostgresRepo) NewID() (int, error) {
	r.Lock()
	defer r.Unlock()
	r.lastID++
	return r.lastID, nil
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
