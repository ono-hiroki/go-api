// Package di は依存性注入のワイヤリングを提供する。
package di

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Container は依存関係のコンテナ。
type Container struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewContainer はコンテナを生成する。
func NewContainer(pool *pgxpool.Pool, logger *slog.Logger) *Container {
	return &Container{pool: pool, logger: logger}
}

// Logger はロガーを返す。
func (c *Container) Logger() *slog.Logger {
	return c.logger
}

// Pool はDBコネクションプールを返す。
func (c *Container) Pool() *pgxpool.Pool {
	return c.pool
}
