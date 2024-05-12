package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/P1coFly/gRPCnmap/internal/app"
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

	application := app.New(log, cfg.GRPC.Port, cfg.GRPC.Timeout)

	go application.GRPCServer.MustRun()

	//GracefulStop
	//waiting for the signal to stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop

	log.Debug("stopping server", slog.String("signal", sig.String()))
	application.GRPCServer.Stop()
	log.Info("server stopped")
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
