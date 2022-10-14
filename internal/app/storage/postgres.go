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

const (
	createUrls = `CREATE TABLE IF NOT EXISTS urls (
            	short varchar(255) PRIMARY KEY,
                original varchar(255),
    			userid varchar(64),
    			UNIQUE(original)
                )`
	addQuery = `
	INSERT INTO urls (short, original, userid)
	VALUES ($1, $2, $3)
	ON CONFLICT DO NOTHING
	RETURNING short`
	countURL       = `SELECT COUNT(*) FROM urls`
	getQuery       = `SELECT original FROM urls WHERE short = $1`
	getByUserQuery = `SELECT * FROM urls WHERE userid = $1`
	getShort       = `SELECT short FROM urls WHERE original = $1`
	drop           = `DROP TABLE urls`
)

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

func (r *PostgresRepo) Get(id string) (string, error) {
	var url string
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := r.conn.QueryRow(ctx, getQuery, id).Scan(&url)
	if err != nil {
		return "", fmt.Errorf("invalid id: %v", id)
	}
	return url, nil
}

func (r *PostgresRepo) GetByUser(userID string) ([]*model.URL, error) {
	var url model.URL
	urls := make([]*model.URL, 0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

func (r *PostgresRepo) Add(url *model.URL) (bool, error) {
	var duplicates bool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := r.conn.QueryRow(ctx, addQuery, url.ShortURL, url.LongURL, url.UserID).Scan(&url.ShortURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			duplicates = true
			err = r.conn.QueryRow(ctx, getShort, url.LongURL).Scan(&url.ShortURL)
		}
		if err != nil {
			return false, err
		}
	}
	return duplicates, nil
}

func (r *PostgresRepo) AddBatch(urls map[string]*model.URL) (bool, error) {
	var duplicates bool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return false, err
	}
	for _, v := range urls {
		err = tx.QueryRow(ctx, addQuery, v.ShortURL, v.LongURL, v.UserID).Scan(&v.ShortURL)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				duplicates = true
				err = tx.QueryRow(ctx, getShort, v.LongURL).Scan(&v.ShortURL)
			}
			if err != nil {
				if err = tx.Rollback(ctx); err != nil {
					return false, err
				}
			}
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return false, err
	}
	return duplicates, nil
}

func (r *PostgresRepo) NewID() (int, error) {
	r.Lock()
	defer r.Unlock()
	r.lastID++
	return r.lastID, nil
}

func (r *PostgresRepo) Ping() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return r.conn.Ping(ctx)
}

func (r *PostgresRepo) DeleteRepo() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := r.conn.Exec(ctx, drop)
	if err != nil {
		return err
	}
	return nil
}
