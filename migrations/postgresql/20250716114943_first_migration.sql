-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS rates(
  id SERIAL PRIMARY KEY,
  market VARCHAR NOT NULL DEFAULT '',
  ask_price DECIMAL(19, 4),
  bid_price DECIMAL(19, 4),
  ts BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)
);

CREATE INDEX IF NOT EXISTS idx_rates_market ON rates(market);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_rates_market;
DROP TABLE IF EXISTS rates;
-- +goose StatementEnd
