package app

import (
	"context"
	"log"

	pb "google.golang.org/grpc/health/grpc_health_v1"
)

type Health struct{}

func New() *Health {
	return &Health{}
}
func (h *Health) Check(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	log.Printf("checking............%s", in.Service)
	var s pb.HealthCheckResponse_ServingStatus = 1
	if in.Service != "service-a" {
		return &pb.HealthCheckResponse{
			Status: 3,
		}, nil
	}
	return &pb.HealthCheckResponse{
		Status: s,
	}, nil
}
func (h *Health) Watch(in *pb.HealthCheckRequest, w pb.Health_WatchServer) error {
	log.Printf("watching............%s", in.Service)
	var s pb.HealthCheckResponse_ServingStatus = 1
	r := &pb.HealthCheckResponse{
		Status: s,
	}
	for {
		w.Send(r)
	}
	return nil
}
