// Package garantex is a package that provides a client for the Garantex API.
package garantex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KVSH-user/ExchangeRateService/internal/config"
)

type Client struct {
	ctx            context.Context
	httpClient     *http.Client
	BaseURL        string
	RequestTimeout time.Duration
}

func NewClient(ctx context.Context, cfg *config.Config) *Client {
	client := &Client{
		ctx: ctx,
		httpClient: &http.Client{
			Timeout: cfg.GarantexClient.Timeout,
		},
		BaseURL: cfg.GarantexClient.BaseURL,
	}

	return client
}

func (cl *Client) GetExchangeRate(_ context.Context, marketID string) (*Response, error) {
	endpoint := fmt.Sprintf("/api/v2/depth?market=%s", marketID)

	var resp Response
	err := cl.doGet(endpoint, &resp)
	if err != nil {
		return nil, fmt.Errorf("doGet: %w", err)
	}

	return &resp, nil
}

func (cl *Client) doGet(endpoint string, v interface{}) error {
	req, err := http.NewRequestWithContext(cl.ctx, http.MethodGet, cl.BaseURL+endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	if err = cl.checkStatusCode(resp); err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

func (cl *Client) checkStatusCode(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnprocessableEntity:
		return ErrInvalidMarketID
	default:
		return fmt.Errorf("status: %d, body: %s", resp.StatusCode, resp.Status)
	}
}
