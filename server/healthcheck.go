package server

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

// HealthChecker ...
type HealthChecker struct{}

// Check ...
func (s *HealthChecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch ...
func (s *HealthChecker) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

// NewHealthChecker ...
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{}
}
