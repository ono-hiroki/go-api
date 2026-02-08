// Package config はアプリケーション設定を提供する。
package config

import (
	"os"
	"strconv"
	"time"
)

// Config はアプリケーション全体の設定。
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig はHTTPサーバーの設定。
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig はデータベース接続の設定。
type DatabaseConfig struct {
	URL             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
}

// Load は環境変数から設定を読み込む。
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "postgres://app:password@localhost:5432/go_api?sslmode=disable"),
			MaxConns:        getInt32Env("DATABASE_MAX_CONNS", 10),
			MinConns:        getInt32Env("DATABASE_MIN_CONNS", 2),
			MaxConnLifetime: getDurationEnv("DATABASE_MAX_CONN_LIFETIME", 30*time.Minute),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getInt32Env(key string, defaultValue int32) int32 {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			return int32(i)
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return defaultValue
}
