package grpcview

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/P1coFly/gRPCnmap/internal/grpc/checkvuln"
	"google.golang.org/grpc"
)

type View struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, controller checkvuln.Controller, port int) *View {
	gRPCServer := grpc.NewServer()

	checkvuln.Register(gRPCServer, controller)

	return &View{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// Run runs gRPC server.
func (v *View) Run() error {
	const op = "grpcview.Run"

	log := v.log.With(slog.String("op", op), slog.Int("port", v.port))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", v.port))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", listener.Addr().String()))

	if err := v.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// MustRun runs gRPC server and panics if any error occurs.
func (v *View) MustRun() {
	if err := v.Run(); err != nil {
		panic(err)
	}
}

// Stop stops gRPC server.
func (v *View) Stop() {
	const op = "grpcview.Stop"

	log := v.log.With(slog.String("op", op), slog.Int("port", v.port))

	log.Info("stopping gRPC server")

	v.gRPCServer.GracefulStop()
}
