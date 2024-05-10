package main

import (
	"log/slog"
	"os"

	"github.com/P1coFly/gRPCnmap/internal/config"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.LogLvl)
	log.Info("starting server")

	log.Debug("cfg's fields", "cfg", cfg)
}

func setupLogger(logLVL string) *slog.Logger {
	var log *slog.Logger

	switch logLVL {
	case LevelDebug:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case LevelInfo:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case LevelWarn:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	case LevelError:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	}

	return log

}
