package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
}

type Config struct {
	URL             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func New(ctx context.Context, config Config) (*Pool, error) {
	//Конфигурация для пула
	poolConfig, err := pgxpool.ParseConfig(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}
	if poolConfig.MaxConns > 0 {
		poolConfig.MaxConns = config.MaxConns
	}
	if poolConfig.MinConns > 0 {
		poolConfig.MinConns = config.MinConns
	}
	if poolConfig.MaxConnLifetime > 0 {
		poolConfig.MaxConnLifetime = config.MaxConnLifetime
	}
	if poolConfig.MaxConnIdleTime > 0 {
		poolConfig.MaxConnIdleTime = config.MaxConnIdleTime
	}
	//Создание пула по новому конфигу
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool with config: %w", err)
	}
	//Проверяем подключение
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Pool{pool}, nil
}

func (p *Pool) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
