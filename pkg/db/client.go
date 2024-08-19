package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(ctx context.Context) (*Database, error) {
	pool, err := pgxpool.New(ctx, generateDsn())
	if err != nil {
		return nil, err
	}
	return newDatabase(pool), nil
}

func generateDsn() string {
	dsn, ok := os.LookupEnv("POSTGRES_DB_DSN")
	if !ok {
		panic("No POSTGRES_DB_DSN in env file")
	}
	return dsn
}
