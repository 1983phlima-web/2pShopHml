package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/2pshop/2pshop/internal/catalog/adapters/postgres"
	catalogApp "github.com/2pshop/2pshop/internal/catalog/application"
	catalogHTTP "github.com/2pshop/2pshop/internal/catalog/transport/http"
	"github.com/2pshop/2pshop/internal/checkout/application"
	"github.com/2pshop/2pshop/internal/identity/application"
	identityHTTP "github.com/2pshop/2pshop/internal/identity/transport/http"
	"github.com/2pshop/2pshop/internal/inventory/application"
	"github.com/2pshop/2pshop/internal/orders/application"
	"github.com/2pshop/2pshop/internal/payments/application"
	"github.com/2pshop/2pshop/internal/platform"
	"github.com/2pshop/2pshop/internal/tenancy/adapters/postgres"
	"github.com/2pshop/2pshop/internal/tenancy/application"
	"github.com/2pshop/2pshop/internal/tenancy/transport/http"
	"github.com/2pshop/2pshop/pkg/idempotency"
	"github.com/2pshop/2pshop/pkg/telemetry"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
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

	// Observability bootstrap
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

	// Database
	dbMetrics, _ := telemetry.NewDBMetrics(tel.MeterProvider.Meter("db"))
	db, err := platform.NewDB(context.Background(), cfg.DatabaseURL, dbMetrics)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Redis
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		logger.Fatal("failed to parse redis url", zap.Error(err))
	}
	redisClient := redis.NewClient(opt)
	defer redisClient.Close()

	// Repositories
	tenancyRepo := tenancyPostgres.NewRepository(db)
	catalogRepo := catalogPostgres.NewRepository(db)

	// Services
	tenancyService := tenancyApp.NewService(tenancyRepo)
	catalogService := catalogApp.NewService(catalogRepo, nil)
	inventoryService := inventoryApp.NewService(nil) // TODO: implement repository
	orderService := ordersApp.NewService(nil, nil)   // TODO: implement repository
	paymentService := paymentsApp.NewService(nil)    // TODO: implement providers
	idempotencyManager := idempotency.NewManager(idempotency.NewRedisStore(redisClient, 24*time.Hour))

	checkoutService := checkoutApp.NewService(
		catalogRepo,
		inventoryService,
		orderService,
		paymentService,
		idempotencyManager,
		tel.Tracer("checkout"),
	)

	// HTTP Metrics
	httpMetrics, _ := telemetry.NewHTTPMetrics(tel.MeterProvider.Meter("http"))

	// Router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Idempotency-Key", "Traceparent", "X-Request-ID", "X-Tenant-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health
	health := platform.NewHealthHandler()
	health.Register("database", platform.NewDBHealthChecker(db))
	r.Get("/health", health.Handle)
	r.Get("/ready", health.Handle)

	// Version
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"version":"%s","build_time":"%s","git_commit":"%s"}`+"\n", version, buildTime, gitCommit)
	})

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return telemetry.InstrumentHandler(next, httpMetrics, "api")
		})

		// Tenancy
		tenancyHandler := tenancyHTTP.NewHandler(tenancyService)
		tenancyHandler.Routes(r)

		// Catalog
		catalogHandler := catalogHTTP.NewHandler(catalogService)
		catalogHandler.Routes(r)

		// Checkout
		r.Post("/checkout", func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement checkout handler
			w.WriteHeader(http.StatusNotImplemented)
		})
	})

	// Server
	port := cfg.HTTPPort
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		logger.Info("starting server", zap.String("addr", srv.Addr), zap.String("version", version))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", zap.Error(err))
	}
}
