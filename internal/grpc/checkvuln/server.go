package checkvuln

import (
	"context"

	"github.com/P1coFly/gRPCnmap/pkg/netvuln_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Controller interface {
	CheckVuln(ctx context.Context, targets []string, ports []int32) ([]*netvuln_v1.TargetResult, error)
}

type serverAPI struct {
	netvuln_v1.UnimplementedNetVulnServiceServer
	controller Controller
}

func Register(gRPC *grpc.Server, controller Controller) {
	netvuln_v1.RegisterNetVulnServiceServer(gRPC, &serverAPI{controller: controller})
}

func (s *serverAPI) CheckVuln(ctx context.Context, req *netvuln_v1.CheckVulnRequest) (*netvuln_v1.CheckVulnResponse, error) {
	if req.GetTargets() == nil {
		return nil, status.Error(codes.InvalidArgument, "targets is required")
	}

	if req.GetTcpPort() == nil {
		return nil, status.Error(codes.InvalidArgument, "ports is required")
	}

	targRes, err := s.controller.CheckVuln(ctx, req.GetTargets(), req.GetTcpPort())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &netvuln_v1.CheckVulnResponse{
		Results: targRes,
	}, nil
}
