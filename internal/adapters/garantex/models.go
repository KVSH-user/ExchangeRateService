package garantex

import (
	"fmt"

	newDecimal "github.com/shopspring/decimal"

	"github.com/KVSH-user/ExchangeRateService/internal/models"
)

type Response struct {
	Timestamp int   `json:"timestamp"`
	Asks      []Ask `json:"asks"`
	Bids      []Bid `json:"bids"`
}

type Ask struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

type Bid struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

func (r Response) ToExchangeRateModel() (*models.ExchangeRate, error) {
	var (
		ask, bid newDecimal.Decimal
		err      error
	)

	if len(r.Asks) > 0 {
		ask, err = newDecimal.NewFromString(r.Asks[0].Price)
		if err != nil {
			return nil, fmt.Errorf("could not convert ask price to decimal: %w", err)
		}
	}

	if len(r.Bids) > 0 {
		bid, err = newDecimal.NewFromString(r.Bids[0].Price)
		if err != nil {
			return nil, fmt.Errorf("could not convert bid price to decimal: %w", err)
		}
	}

	return &models.ExchangeRate{
		AskPrice: ask,
		BidPrice: bid,
		TS:       int64(r.Timestamp),
	}, nil
}
