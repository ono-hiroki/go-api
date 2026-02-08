package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	pool, err := database.Connect()
	if err != nil {
		return err
	}
	defer pool.Close()

	container := di.NewContainer(pool, logger)

	h := httpapi.NewRouter(container)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Info("server starting", "addr", ":8080")
	return server.ListenAndServe()
}
