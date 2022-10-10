package storage

import (
	"context"
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	conn *pgxpool.Pool
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

func (r *PostgresRepo) Get(id string) (string, error) {
	return "", nil
}

func (r *PostgresRepo) GetByUser(userID string) ([]*model.URL, error) {
	return nil, nil
}

func (r *PostgresRepo) Add(long, short, userID string) (string, error) {
	return "", nil
}

func (r *PostgresRepo) NewID() (int, error) {
	return 0, nil
}

func (r *PostgresRepo) Ping() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return r.conn.Ping(ctx)
}
