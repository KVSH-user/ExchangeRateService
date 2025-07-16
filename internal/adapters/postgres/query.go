package postgres

import (
	"context"
	"fmt"

	"github.com/KVSH-user/ExchangeRateService/internal/models"
)

// SaveExchangeRate - method for save exchange rate to db
func (s *Store) SaveExchangeRate(ctx context.Context, rate *models.ExchangeRate) error {
	const query = `
		INSERT INTO rates (
			ask_price, bid_price, ts
		) VALUES (
			$1, $2, $3
		) RETURNING id`

	err := s.queryRow(ctx, query, s.Master,
		rate.AskPrice,
		rate.BidPrice,
		rate.TS,
	).Scan(&rate.ID)

	if err != nil {
		return fmt.Errorf("SaveExchangeRate: %w", err)
	}

	return nil
}
