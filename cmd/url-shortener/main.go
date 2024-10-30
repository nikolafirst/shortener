package main

import (
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	// fmt.Println(cfg)

	// init logger: slog

	log := setupLogger(cfg.Env)

	log.Info("URL Shortener started", slog.String("env", cfg.Env))
	log.Debug("Debug message are enabled")

	// init storage: sqlite

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// TODO: init router: chi, "chi render"

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
