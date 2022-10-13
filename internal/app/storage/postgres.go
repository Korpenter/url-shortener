package storage

import (
	"context"
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	conn   *pgxpool.Pool
	lastID int
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
	urls := `CREATE TABLE IF NOT EXISTS urls (
            	short varchar(255) PRIMARY KEY,
                original varchar(255),
    			userid varchar(64)
                )`
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := r.conn.Exec(ctx, urls)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepo) NumberOfURLs() error {
	amount := `SELECT COUNT(*) FROM urls`
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rows, err := r.conn.Query(ctx, amount)
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
	getQuery := `SELECT original FROM urls WHERE short = $1`
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := r.conn.QueryRow(ctx, getQuery, id).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (r *PostgresRepo) GetByUser(userID string) ([]*model.URL, error) {
	var url model.URL
	urls := make([]*model.URL, 0)

	getByUserQuery := `SELECT * FROM urls WHERE userid = $1`
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

func (r *PostgresRepo) Add(url *model.URL) error {
	addQuery := `
	INSERT INTO urls (short, original, userid)
	VALUES ($1, $2, $3)
	RETURNING short`
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := r.conn.Query(ctx, addQuery, url.ShortURL, url.LongURL, url.UserID)
	if err != nil {
		return err
	}
	r.lastID++
	return nil
}

func (r *PostgresRepo) AddBatch(urls []*model.URL) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := r.conn.CopyFrom(
		ctx,
		pgx.Identifier{"people"},
		[]string{"first_name", "last_name", "age"},
		pgx.CopyFromSlice(len(urls), func(i int) ([]any, error) {
			return []any{urls[i].ShortURL, urls[i].LongURL, urls[i].UserID}, nil
		}),
	)
	if err != nil {
		return err
	}
	r.lastID++
	return nil
}

func (r *PostgresRepo) NewID() (int, error) {
	return r.lastID + 1, nil
}

func (r *PostgresRepo) Ping() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return r.conn.Ping(ctx)
}
