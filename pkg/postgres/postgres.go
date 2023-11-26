package postgres

import (
	"context"
	"fmt"
	"test_wb/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectionPG(cfg *config.PG) (*pgxpool.Pool, error) {
	connect := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.PgName)

	return pgxpool.Connect(context.Background(), connect)
}
