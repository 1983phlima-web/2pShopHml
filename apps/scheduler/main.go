package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/2pshop/2pshop/internal/platform"
	"github.com/2pshop/2pshop/pkg/telemetry"
	"go.uber.org/zap"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := platform.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	tel, err := telemetry.New(telemetry.Config{
		ServiceName:    cfg.OTEL.ServiceName,
		ServiceVersion: version,
		Environment:    cfg.OTEL.Environment,
		Region:         cfg.OTEL.Region,
		GitCommit:      gitCommit,
		BuildID:        buildTime,
		OTLPEndpoint:   cfg.OTEL.Endpoint,
		SamplingRate:   cfg.OTEL.SamplingRate,
	}, logger)
	if err != nil {
		logger.Fatal("failed to initialize telemetry", zap.Error(err))
	}
	defer tel.Shutdown(context.Background())

	logger.Info("scheduler started", zap.String("version", version))

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			// TODO: run scheduled tasks
			// - expire inventory reservations
			// - reconciliation
			// - cleanup idempotency keys
			logger.Info("tick")
		case <-quit:
			logger.Info("shutting down scheduler")
			return
		}
	}
}
