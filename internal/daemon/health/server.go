package health

import (
	"context"

	"google.golang.org/grpc/codes"
	hc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type server struct {
}

// Check returns if the service is available.
func (s *server) Check(c context.Context, r *hc.HealthCheckRequest) (*hc.HealthCheckResponse, error) {
	return &hc.HealthCheckResponse{
		Status: hc.HealthCheckResponse_SERVING,
	}, nil
}

// Watch is used to check health using streaming. Not using it for now.
func (s *server) Watch(*hc.HealthCheckRequest, hc.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

// NewServer initializes health server
func NewServer() hc.HealthServer {
	return &server{}
}
