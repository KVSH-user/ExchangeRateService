// Package main
package main

import (
	"context"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KVSH-user/ExchangeRateService/internal/adapters/garantex"
	"github.com/KVSH-user/ExchangeRateService/internal/adapters/postgres"
	"github.com/KVSH-user/ExchangeRateService/internal/app/grpc/exchangerateservice"
	"github.com/KVSH-user/ExchangeRateService/internal/config"
	"github.com/KVSH-user/ExchangeRateService/internal/modules/exchangerate"
	"github.com/KVSH-user/ExchangeRateService/internal/utils"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	log := utils.SetupLogger(cfg.Env)

	storage, err := postgres.NewClient(ctx, log, cfg)
	if err != nil {
		log.Error("Failed to connect to storage", "error", err)
		os.Exit(1)
	}

	garantexClient := garantex.NewClient(ctx, cfg)

	exchangeRateModule := exchangerate.New(log, storage, garantexClient)

	server := exchangerateservice.NewServer(log, cfg.GRPC.Port, exchangeRateModule)

	errChan := make(chan error, 1)

	go func() {
		if err := server.Start(ctx); err != nil {
			errChan <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Error("Server error", "error", err)
	case sig := <-sigChan:
		log.Info("Received signal", "signal", sig.String())
	}

	log.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Stop(shutdownCtx); err != nil {
		log.Error("Error during shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("Server stopped gracefully")
}
