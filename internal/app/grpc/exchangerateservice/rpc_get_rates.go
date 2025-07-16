package exchangerateservice

import (
	"context"

	"google.golang.org/genproto/googleapis/type/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/KVSH-user/ExchangeRateService/pkg/pb/exchangerateservice"
)

func (s *ExchangeRateService) GetRates(ctx context.Context, req *pb.GetRatesRequest) (*pb.GetRatesResponse, error) {
	if err := validateGetRatesReq(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rate, err := s.exchangeRateModule.GetExchangeRate(ctx, req.GetMarket())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch rates: %v", err)
	}

	return &pb.GetRatesResponse{
		Ts: rate.TS,
		AskPrice: &decimal.Decimal{
			Value: rate.AskPrice.String(),
		},
		BidPrice: &decimal.Decimal{
			Value: rate.BidPrice.String(),
		},
	}, nil
}

func validateGetRatesReq(req *pb.GetRatesRequest) error {
	switch {
	case req.GetMarket() == "":
		return status.Errorf(codes.InvalidArgument, "market is required")
	default:
		return nil
	}
}
