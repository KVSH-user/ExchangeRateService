package models

import "github.com/shopspring/decimal"

type ExchangeRate struct {
	ID       int64           `json:"id"`
	AskPrice decimal.Decimal `json:"ask_price"`
	BidPrice decimal.Decimal `json:"bid_price"`
	TS       int64           `json:"ts"`
}
