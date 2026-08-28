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
	"go.opentelemetry.io/otel/attribute"
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
		GroupID:     "2pshop-workers",
		Topic:       "2pshop.events",
		MinBytes:    10e3,
		MaxBytes:    10e6,
		MaxWait:     time.Second,
		StartOffset: kafka.FirstOffset,
	})
	defer reader.Close()

	logger.Info("worker started", zap.String("version", version), zap.String("brokers", cfg.KafkaBrokers))

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

			_, span := tel.Tracer("worker").Start(ctx, "worker.process")
			span.SetAttributes(
				attribute.String("kafka.topic", msg.Topic),
				attribute.String("kafka.partition", fmt.Sprintf("%d", msg.Partition)),
				attribute.Int64("kafka.offset", msg.Offset),
				attribute.String("event.type", string(msg.Key)),
			)

			logger.Info("processing message",
				zap.String("topic", msg.Topic),
				zap.Int("partition", msg.Partition),
				zap.Int64("offset", msg.Offset),
				zap.String("key", string(msg.Key)),
			)

			// TODO: route to specific handlers based on event type
			// TODO: implement idempotent processing
			// TODO: implement DLQ for failed messages

			span.End()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down worker")
	cancel()
}
