package db

import (
	"context"
	"fmt"

	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/config"
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
	cfg := config.MustLoad()
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		cfg.DatabaseData.Host,
		cfg.DatabaseData.Port,
		cfg.DatabaseData.User,
		cfg.DatabaseData.Password,
		cfg.DatabaseData.DBName,
	)
}
