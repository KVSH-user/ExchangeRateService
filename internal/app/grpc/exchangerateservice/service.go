package exchangerateservice

import (
	"context"
	"log/slog"

	"github.com/KVSH-user/ExchangeRateService/internal/models"

	pb "github.com/KVSH-user/ExchangeRateService/pkg/pb/exchangerateservice"
)

type ExchangeRateService struct {
	pb.UnimplementedExchangeRateServiceServer
	logger             *slog.Logger
	exchangeRateModule ExchangeRateModule
}

type ExchangeRateModule interface {
	GetExchangeRate(ctx context.Context, market string) (*models.ExchangeRate, error)
}

func NewExchangeRateService(logger *slog.Logger, exchangeRateModule ExchangeRateModule) *ExchangeRateService {
	return &ExchangeRateService{
		logger:             logger,
		exchangeRateModule: exchangeRateModule,
	}
}
