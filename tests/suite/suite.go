package suite

import (
	"context"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/P1coFly/gRPCnmap/internal/config"
	"github.com/P1coFly/gRPCnmap/pkg/netvuln_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T                                           // Потребуется для вызова методов *testing.T внутри Suite
	Cfg                  *config.Config                  // Конфигурация приложения
	NetVulnServiceClient netvuln_v1.NetVulnServiceClient // Клиент для взаимодействия с gRPC-сервером
}

const (
	grpcHost = "localhost"
)

// New creates new test suite.
//
// TODO: for pipeline tests we need to wait for app is ready
func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath(configPath())

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(),
		net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials())) // Используем insecure-коннект для тестов
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:                    t,
		Cfg:                  cfg,
		NetVulnServiceClient: netvuln_v1.NewNetVulnServiceClient(cc),
	}
}

func configPath() string {
	const key = "CONFIG_PATH"

	if v := os.Getenv(key); v != "" {
		return v
	}

	return "../config/config.yaml"
}
