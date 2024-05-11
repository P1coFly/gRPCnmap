package tests

import (
	"testing"

	"github.com/P1coFly/gRPCnmap/pkg/netvuln_v1"
	"github.com/P1coFly/gRPCnmap/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ip = "127.0.0.1"
)

func TestCheckVuln_Success(t *testing.T) {
	ctx, st := suite.New(t)

	resp, err := st.NetVulnServiceClient.CheckVuln(ctx, &netvuln_v1.CheckVulnRequest{
		Targets: []string{ip},
		TcpPort: []int32{int32(st.Cfg.GRPC.Port)},
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.GetResults())

	assert.Greater(t, len(resp.GetResults()), 0)

	respIP := resp.GetResults()[0].Target
	assert.Equal(t, respIP, ip)

	assert.Greater(t, len(resp.GetResults()[0].Services), 0)

	port := resp.GetResults()[0].Services[0].TcpPort
	assert.Equal(t, port, int32(st.Cfg.GRPC.Port))
}

func TestCheckVuln_WithoutTargets(t *testing.T) {
	ctx, st := suite.New(t)

	resp, err := st.NetVulnServiceClient.CheckVuln(ctx, &netvuln_v1.CheckVulnRequest{
		Targets: []string{},
		TcpPort: []int32{int32(st.Cfg.GRPC.Port)},
	})

	require.Error(t, err)
	assert.Empty(t, resp.GetResults())
	assert.ErrorContains(t, err, "targets is required")
}

func TestCheckVuln_WithoutPorts(t *testing.T) {
	ctx, st := suite.New(t)

	resp, err := st.NetVulnServiceClient.CheckVuln(ctx, &netvuln_v1.CheckVulnRequest{
		Targets: []string{ip},
		TcpPort: []int32{},
	})

	require.Error(t, err)
	assert.Empty(t, resp.GetResults())
	assert.ErrorContains(t, err, "ports is required")
}
