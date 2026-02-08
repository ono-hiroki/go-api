package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"go-api/internal/config"
	"go-api/internal/di"
	"go-api/internal/infrastructure/database"
	httpapi "go-api/internal/presentation/http"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	pool, err := database.Connect(context.Background(), cfg.Database)
	if err != nil {
		return err
	}
	defer pool.Close()

	container := di.NewContainer(pool, logger)

	h := httpapi.NewRouter(container)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      h,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	logger.Info("server starting", "addr", server.Addr)
	return server.ListenAndServe()
}
