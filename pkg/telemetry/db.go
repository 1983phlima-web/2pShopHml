package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type DBMetrics struct {
	QueryDuration    metric.Float64Histogram
	PoolSaturation   metric.Float64ObservableGauge
	ConnectionErrors metric.Int64Counter
}

func NewDBMetrics(meter metric.Meter) (*DBMetrics, error) {
	dur, err := meter.Float64Histogram("db.query.duration",
		metric.WithDescription("Database query duration"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	errs, err := meter.Int64Counter("db.connection.errors",
		metric.WithDescription("Database connection errors"),
	)
	if err != nil {
		return nil, err
	}

	return &DBMetrics{
		QueryDuration:    dur,
		ConnectionErrors: errs,
	}, nil
}

func RecordQuery(ctx context.Context, metrics *DBMetrics, operation string, duration time.Duration, err error) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetAttributes(
			attribute.String("db.operation", operation),
			attribute.Float64("db.query.duration_ms", float64(duration.Milliseconds())),
		)
		if err != nil {
			span.RecordError(err)
		}
	}

	if metrics != nil {
		attrs := []attribute.KeyValue{
			attribute.String("db.operation", operation),
		}
		metrics.QueryDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
		if err != nil {
			metrics.ConnectionErrors.Add(ctx, 1, metric.WithAttributes(
				attribute.String("error.type", err.Error()),
			))
		}
	}
}
