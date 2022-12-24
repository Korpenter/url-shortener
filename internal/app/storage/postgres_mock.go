package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
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
	mockDeleteQuery    = `UPDATE urls_test SET deleted=TRUE WHERE (short, userid) IN `
	mockGetQuery       = `SELECT original, deleted FROM urls_test WHERE short = $1`
	mockGetByUserQuery = `SELECT * FROM urls_test WHERE userid = $1`
	mockGetShort       = `SELECT short FROM urls_test WHERE original = $1`
	mockDrop           = `DROP TABLE urls_test`
)

type PostgresMockRepo struct {
	conn *pgxpool.Pool
}

func NewPostgresMockRepo(connString string) (*PostgresMockRepo, error) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	url1 := &model.URL{ShortURL: "1", LongURL: "https://github.com/Mldlr/url-shortener/internal/app/utils/encoders"}
	url2 := &model.URL{ShortURL: "2", LongURL: "https://yandex.ru/"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = conn.Exec(ctx, mockCreateUrls)
	if err != nil {
		return nil, err
	}
	repo := &PostgresMockRepo{conn: conn}
	repo.Add(context.Background(), url1)
	repo.Add(context.Background(), url2)
	return repo, nil
}

func (r *PostgresMockRepo) Get(ctx context.Context, id string) (*model.URL, error) {
	var url model.URL
	err := r.conn.QueryRow(ctx, mockGetQuery, id).Scan(&url.LongURL, &url.Deleted)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", id)
	}
	return &url, nil
}

func (r *PostgresMockRepo) GetByUser(ctx context.Context, userID string) ([]*model.URL, error) {
	var url model.URL
	urls := make([]*model.URL, 0)
	rows, err := r.conn.Query(ctx, mockGetByUserQuery, userID)
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

func (r *PostgresMockRepo) Add(ctx context.Context, url *model.URL) (bool, error) {
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
	err = tx.QueryRow(ctx, mockAddQuery, url.ShortURL, url.LongURL, url.UserID).Scan(&url.ShortURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			duplicates = true
			err = tx.QueryRow(ctx, mockGetShort, url.LongURL).Scan(&url.ShortURL)
		} else if err != nil {
			return false, err
		}
	}
	return duplicates, nil
}

func (r *PostgresMockRepo) AddBatch(ctx context.Context, urls map[string]*model.URL) (bool, error) {
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
		err = tx.QueryRow(ctx, mockAddQuery, v.ShortURL, v.LongURL, v.UserID).Scan(&v.ShortURL)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				duplicates = true
				err = tx.QueryRow(ctx, mockGetShort, v.LongURL).Scan(&v.ShortURL)
			}
		}
	}
	return duplicates, nil
}

func (r *PostgresMockRepo) NewID(url string) (string, error) {
	return encoders.ToRBase62(url), nil
}

func (r *PostgresMockRepo) DeleteURLs(deleteURLs []*model.DeleteURLItem) (int, error) {
	ctx := context.Background()
	tx, err := r.conn.Begin(ctx)
	var n int
	if err != nil {
		return n, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	var args string
	for _, v := range deleteURLs {
		args += fmt.Sprintf("('%s','%s'),", v.ShortURL, v.UserID)
	}
	args = args[:len(args)-1]
	query := mockDeleteQuery + `(` + args + `)`
	res, err := tx.Exec(ctx, query)
	if err != nil {
		return n, err
	}
	n = int(res.RowsAffected())
	return n, err
}

func (r *PostgresMockRepo) Ping(ctx context.Context) error {
	return r.conn.Ping(ctx)
}

func (r *PostgresMockRepo) DeleteRepo(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, mockDrop)
	if err != nil {
		return err
	}
	return nil
}
