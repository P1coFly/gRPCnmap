package app

import (
	"log/slog"
	"time"

	grpcview "github.com/P1coFly/gRPCnmap/internal/app/grpc"
	"github.com/P1coFly/gRPCnmap/internal/controller/nmapservices"
)

type App struct {
	GRPCServer *grpcview.View
}

func New(
	log *slog.Logger,
	grpcPort int,
	timeout time.Duration,
) *App {

	controller := nmapservices.New(log, timeout)
	grpcApp := grpcview.New(log, controller, grpcPort)
	return &App{
		GRPCServer: grpcApp,
	}

}
