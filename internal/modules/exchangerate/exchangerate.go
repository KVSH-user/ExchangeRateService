package exchangerate

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/KVSH-user/ExchangeRateService/internal/adapters/garantex"
	"github.com/KVSH-user/ExchangeRateService/internal/models"
)

type RateStorage interface {
	SaveExchangeRate(ctx context.Context, rate *models.ExchangeRate) error
}

type GarantexClient interface {
	GetExchangeRate(_ context.Context, marketID string) (*garantex.Response, error)
}

type Module struct {
	log            *slog.Logger
	rateStorage    RateStorage
	garantexClient GarantexClient
}

func New(log *slog.Logger, rateStorage RateStorage, garantexClient GarantexClient) *Module {
	return &Module{
		log:            log,
		rateStorage:    rateStorage,
		garantexClient: garantexClient,
	}
}

func (m *Module) GetExchangeRate(ctx context.Context, market string) (*models.ExchangeRate, error) {
	resp, err := m.garantexClient.GetExchangeRate(ctx, market)
	if err != nil {
		m.log.ErrorContext(ctx, "failed to fetch exchange rate", "error", err)

		return nil, fmt.Errorf("could not get exchange rate: %w", err)
	}

	rate, err := resp.ToExchangeRateModel()
	if err != nil {
		m.log.ErrorContext(ctx, "failed to convert exchange rate to model", "error", err)

		return nil, fmt.Errorf("could not convert exchange rate to model: %w", err)
	}

	err = m.rateStorage.SaveExchangeRate(ctx, rate)
	if err != nil {
		m.log.ErrorContext(ctx, "failed to save exchange rate", "error", err)

		return nil, fmt.Errorf("could not save exchange rate: %w", err)
	}

	return rate, nil
}
