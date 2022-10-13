package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresMockRepo struct {
	conn   *pgxpool.Pool
	lastID int
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
	urls := `CREATE TABLE IF NOT EXISTS urls_test (
            	short varchar(255) PRIMARY KEY,
                original varchar(255),
    			userid varchar(64),
    			UNIQUE(original)
                )`
	url1 := &model.URL{ShortURL: "1", LongURL: "https://github.com/Mldlr/url-shortener/internal/app/utils/encoders"}
	url2 := &model.URL{ShortURL: "2", LongURL: "https://yandex.ru/"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = conn.Exec(ctx, urls)
	if err != nil {
		return nil, err
	}
	repo := &PostgresMockRepo{conn: conn}
	repo.Add(url1)
	repo.Add(url2)
	repo.lastID = 2
	return repo, nil
}

func (r *PostgresMockRepo) Get(id string) (string, error) {
	var url string
	getQuery := `SELECT original FROM urls_test WHERE short = $1`
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := r.conn.QueryRow(ctx, getQuery, id).Scan(&url)
	if err != nil {
		return "", fmt.Errorf("invalid id: %v", id)
	}
	return url, nil
}

func (r *PostgresMockRepo) GetByUser(userID string) ([]*model.URL, error) {
	var url model.URL
	urls := make([]*model.URL, 0)

	getByUserQuery := `SELECT * FROM urls_test WHERE userid = $1`
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

func (r *PostgresMockRepo) Add(url *model.URL) error {
	addQuery := `
	INSERT INTO urls_test (short, original, userid)
	VALUES ($1, $2, $3)
	ON CONFLICT DO NOTHING
	RETURNING short`
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := r.conn.QueryRow(ctx, addQuery, url.ShortURL, url.LongURL, url.UserID).Scan(&url.ShortURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = r.conn.QueryRow(ctx, `SELECT short FROM urls_test WHERE original = $1`, url.LongURL).Scan(&url.ShortURL)
		} else {
			return err
		}
	}
	return nil
}

func (r *PostgresMockRepo) AddBatch(urls []model.URL) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := r.conn.CopyFrom(
		ctx,
		pgx.Identifier{"urls_test"},
		[]string{"short", "original", "userid"},
		pgx.CopyFromSlice(len(urls), func(i int) ([]any, error) {
			fmt.Println(i, urls[i])
			return []any{urls[i].ShortURL, urls[i].LongURL, urls[i].UserID}, nil
		}),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresMockRepo) NewID() (int, error) {
	r.lastID++
	return r.lastID, nil
}

func (r *PostgresMockRepo) Ping() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return r.conn.Ping(ctx)
}

func (r *PostgresMockRepo) Delete() error {
	urls := `DROP TABLE urls_test`
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := r.conn.Exec(ctx, urls)
	if err != nil {
		return err
	}
	return nil
}
