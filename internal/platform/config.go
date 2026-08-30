package platform

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv       string
	HTTPPort     string
	DatabaseURL  string
	RedisURL     string
	KafkaBrokers string
	JWTSecret    string
	OTEL         OTELConfig
}

type OTELConfig struct {
	Endpoint       string
	ServiceName    string
	ServiceVersion string
	Environment    string
	Region         string
	SamplingRate   float64
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		AppEnv:       getEnv("APP_ENV", "development"),
		HTTPPort:     getEnv("HTTP_PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		RedisURL:     getEnv("REDIS_URL", ""),
		KafkaBrokers: getEnv("KAFKA_BROKERS", ""),
		JWTSecret:    getEnv("JWT_SECRET", "2pshop-hml-default-secret-change-me"),
		OTEL: OTELConfig{
			Endpoint:       stripScheme(getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")),
			ServiceName:    getEnv("OTEL_SERVICE_NAME", "2pshop"),
			ServiceVersion: getEnv("OTEL_SERVICE_VERSION", "dev"),
			Environment:    getEnv("APP_ENV", "development"),
			Region:         getEnv("RAILWAY_REGION", "us-east-1"),
			SamplingRate:   getEnvFloat("OTEL_SAMPLING_RATE", 1.0),
		},
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// stripScheme removes a leading http:// or https:// from an endpoint,
// since OTLP gRPC exporters expect a bare host:port.
func stripScheme(endpoint string) string {
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	return endpoint
}

func getEnvFloat(key string, defaultVal float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return defaultVal
	}
	return f
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return defaultVal
	}
	return d
}
