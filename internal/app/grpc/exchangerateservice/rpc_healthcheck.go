package exchangerateservice

import (
	"context"

	pb "github.com/KVSH-user/ExchangeRateService/pkg/pb/exchangerateservice"
)

func (s *ExchangeRateService) HealthCheck(_ context.Context, _ *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: "OK",
	}, nil
}
