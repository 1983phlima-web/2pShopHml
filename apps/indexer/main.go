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
	"github.com/segmentio/kafka-go"
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

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{cfg.KafkaBrokers},
		GroupID:     "2pshop-indexers",
		Topic:       "2pshop.events",
		MinBytes:    10e3,
		MaxBytes:    10e6,
		MaxWait:     time.Second,
		StartOffset: kafka.FirstOffset,
	})
	defer reader.Close()

	logger.Info("indexer started", zap.String("version", version))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				logger.Error("kafka read error", zap.Error(err))
				continue
			}
			// TODO: index to OpenSearch
			logger.Info("indexing event", zap.String("key", string(msg.Key)))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down indexer")
	cancel()
}
