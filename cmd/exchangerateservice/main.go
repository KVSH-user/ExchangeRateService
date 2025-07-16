// Package main
package main

import (
	"context"

	"github.com/KVSH-user/ExchangeRateService/internal/config"
	"github.com/KVSH-user/ExchangeRateService/internal/utils"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	log := utils.SetupLogger(cfg.Env)

	_, _ = log, ctx
}
