package telemetry

import (
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type HTTPMetrics struct {
	RequestDuration metric.Float64Histogram
	ActiveRequests  metric.Int64UpDownCounter
	RequestBodySize metric.Int64Histogram
	ResponseBodySize metric.Int64Histogram
}

func NewHTTPMetrics(meter metric.Meter) (*HTTPMetrics, error) {
	dur, err := meter.Float64Histogram("http.server.request.duration",
		metric.WithDescription("Duration of HTTP requests"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	active, err := meter.Int64UpDownCounter("http.server.active_requests",
		metric.WithDescription("Number of active HTTP requests"),
	)
	if err != nil {
		return nil, err
	}

	reqSize, err := meter.Int64Histogram("http.server.request.body.size",
		metric.WithDescription("Size of HTTP request bodies"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, err
	}

	respSize, err := meter.Int64Histogram("http.server.response.body.size",
		metric.WithDescription("Size of HTTP response bodies"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, err
	}

	return &HTTPMetrics{
		RequestDuration:  dur,
		ActiveRequests:   active,
		RequestBodySize:  reqSize,
		ResponseBodySize: respSize,
	}, nil
}

func InstrumentHandler(next http.Handler, metrics *HTTPMetrics, operation string) http.Handler {
	return otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		metrics.ActiveRequests.Add(r.Context(), 1)
		defer metrics.ActiveRequests.Add(r.Context(), -1)

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()
		attrs := []attribute.KeyValue{
			attribute.String("http.request.method", r.Method),
			attribute.String("http.route", r.URL.Path),
			attribute.String("server.address", r.Host),
		}

		metrics.RequestDuration.Record(r.Context(), duration, metric.WithAttributes(attrs...))
	}), operation)
}
