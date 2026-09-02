package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	checkoutApp "github.com/2pshop/2pshop/internal/checkout/application"
	checkoutHTTP "github.com/2pshop/2pshop/internal/checkout/transport/http"

	analyticsPG "github.com/2pshop/2pshop/internal/analytics/adapters/postgres"
	analyticsApp "github.com/2pshop/2pshop/internal/analytics/application"
	analyticsHTTP "github.com/2pshop/2pshop/internal/analytics/transport/http"

	favoritesPG "github.com/2pshop/2pshop/internal/favorites/adapters/postgres"
	favoritesApp "github.com/2pshop/2pshop/internal/favorites/application"
	favoritesHTTP "github.com/2pshop/2pshop/internal/favorites/transport/http"

	loyaltyPG "github.com/2pshop/2pshop/internal/loyalty/adapters/postgres"
	loyaltyApp "github.com/2pshop/2pshop/internal/loyalty/application"
	loyaltyHTTP "github.com/2pshop/2pshop/internal/loyalty/transport/http"

	questionsPG "github.com/2pshop/2pshop/internal/questions/adapters/postgres"
	questionsApp "github.com/2pshop/2pshop/internal/questions/application"
	questionsHTTP "github.com/2pshop/2pshop/internal/questions/transport/http"

	exchangesPG "github.com/2pshop/2pshop/internal/exchanges/adapters/postgres"
	exchangesApp "github.com/2pshop/2pshop/internal/exchanges/application"
	exchangesHTTP "github.com/2pshop/2pshop/internal/exchanges/transport/http"

	settingsPG "github.com/2pshop/2pshop/internal/settings/adapters/postgres"
	settingsApp "github.com/2pshop/2pshop/internal/settings/application"
	settingsHTTP "github.com/2pshop/2pshop/internal/settings/transport/http"

	catalogPG "github.com/2pshop/2pshop/internal/catalog/adapters/postgres"
	catalogApp "github.com/2pshop/2pshop/internal/catalog/application"
	catalogHTTP "github.com/2pshop/2pshop/internal/catalog/transport/http"

	identityJWT "github.com/2pshop/2pshop/internal/identity/adapters/jwt"
	identityPG "github.com/2pshop/2pshop/internal/identity/adapters/postgres"
	identityApp "github.com/2pshop/2pshop/internal/identity/application"
	identityHTTP "github.com/2pshop/2pshop/internal/identity/transport/http"

	inventoryPG "github.com/2pshop/2pshop/internal/inventory/adapters/postgres"
	inventoryApp "github.com/2pshop/2pshop/internal/inventory/application"

	ordersPG "github.com/2pshop/2pshop/internal/orders/adapters/postgres"
	ordersApp "github.com/2pshop/2pshop/internal/orders/application"
	ordersHTTP "github.com/2pshop/2pshop/internal/orders/transport/http"

	paymentsMock "github.com/2pshop/2pshop/internal/payments/adapters/mock"
	paymentsApp "github.com/2pshop/2pshop/internal/payments/application"
	paymentsDomain "github.com/2pshop/2pshop/internal/payments/domain"

	reviewsPG "github.com/2pshop/2pshop/internal/reviews/adapters/postgres"
	reviewsApp "github.com/2pshop/2pshop/internal/reviews/application"
	reviewsHTTP "github.com/2pshop/2pshop/internal/reviews/transport/http"

	tenancyPG "github.com/2pshop/2pshop/internal/tenancy/adapters/postgres"
	tenancyApp "github.com/2pshop/2pshop/internal/tenancy/application"
	tenancyHTTP "github.com/2pshop/2pshop/internal/tenancy/transport/http"

	"github.com/2pshop/2pshop/internal/platform"
	"github.com/2pshop/2pshop/internal/platform/httpmw"
	"github.com/2pshop/2pshop/pkg/idempotency"
	"github.com/2pshop/2pshop/pkg/telemetry"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
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

	// Migrations run automatically on every boot (idempotent).
	if err := platform.RunMigrations(context.Background(), db); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	// Redis (idempotency store for checkout)
	var idempotencyManager *idempotency.Manager
	if cfg.RedisURL != "" {
		opt, err := redis.ParseURL(cfg.RedisURL)
		if err != nil {
			logger.Fatal("failed to parse redis url", zap.Error(err))
		}
		redisClient := redis.NewClient(opt)
		defer redisClient.Close()
		idempotencyManager = idempotency.NewManager(idempotency.NewRedisStore(redisClient, 24*time.Hour))
	}

	// Repositories
	tenancyRepo := tenancyPG.NewRepository(db)
	identityRepo := identityPG.NewRepository(db)
	catalogRepo := catalogPG.NewRepository(db)
	reviewsRepo := reviewsPG.NewRepository(db)
	inventoryRepo := inventoryPG.NewRepository(db)
	ordersRepo := ordersPG.NewRepository(db)
	settingsRepo := settingsPG.NewRepository(db)
	analyticsRepo := analyticsPG.NewRepository(db, version)
	favoritesRepo := favoritesPG.NewRepository(db)
	loyaltyRepo := loyaltyPG.NewRepository(db)
	questionsRepo := questionsPG.NewRepository(db)
	exchangesRepo := exchangesPG.NewRepository(db)

	// Services
	tenancyService := tenancyApp.NewService(tenancyRepo)
	tokenService := identityJWT.NewTokenService(cfg.JWTSecret, 24*time.Hour)
	identityService := identityApp.NewService(identityRepo, tokenService)
	catalogService := catalogApp.NewService(catalogRepo, nil)
	reviewsService := reviewsApp.NewService(reviewsRepo)
	inventoryService := inventoryApp.NewService(inventoryRepo)
	ordersService := ordersApp.NewService(ordersRepo, nil)
	settingsService := settingsApp.NewService(settingsRepo)
	analyticsService := analyticsApp.NewService(analyticsRepo)
	favoritesService := favoritesApp.NewService(favoritesRepo)
	loyaltyService := loyaltyApp.NewService(loyaltyRepo)
	questionsService := questionsApp.NewService(questionsRepo)
	exchangesService := exchangesApp.NewService(exchangesRepo)
	paymentsService := paymentsApp.NewService(map[string]paymentsDomain.Provider{
		"stripe": paymentsMock.New(), // sandbox provider for HML; swap for a real Stripe adapter in production.
	})
	checkoutService := checkoutApp.NewService(
		catalogRepo,
		inventoryService,
		ordersService,
		paymentsService,
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

	tenancyHandler := tenancyHTTP.NewHandler(tenancyService)
	identityHandler := identityHTTP.NewHandler(identityService)
	catalogHandler := catalogHTTP.NewHandler(catalogService)
	reviewsHandler := reviewsHTTP.NewHandler(reviewsService)
	ordersHandler := ordersHTTP.NewHandler(ordersService)
	checkoutHandler := checkoutHTTP.NewHandler(checkoutService)
	settingsHandler := settingsHTTP.NewHandler(settingsService)
	analyticsHandler := analyticsHTTP.NewHandler(analyticsService)
	favoritesHandler := favoritesHTTP.NewHandler(favoritesService)
	loyaltyHandler := loyaltyHTTP.NewHandler(loyaltyService)
	questionsHandler := questionsHTTP.NewHandler(questionsService)
	exchangesHandler := exchangesHTTP.NewHandler(exchangesService)

	const (
		roleSeller      = "SELLER"
		roleBuyer       = "BUYER"
		roleSystemAdmin = "SYSTEM_ADMIN"
		roleGlobalAdmin = "GLOBAL_ADMIN"
	)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return telemetry.InstrumentHandler(next, httpMetrics, "api")
		})

		// Tenant management is not tenant-scoped itself.
		tenancyHandler.Routes(r)

		// Everything below requires a resolved tenant (X-Tenant-ID header).
		r.Group(func(r chi.Router) {
			r.Use(httpmw.TenantResolver(tenancyService))

			// Public within the tenant: auth, browsing catalog & reviews, active theme.
			identityHandler.Routes(r)
			catalogHandler.PublicRoutes(r)
			reviewsHandler.PublicRoutes(r)
			questionsHandler.PublicRoutes(r)
			settingsHandler.PublicRoutes(r)

			// Authenticated routes.
			r.Group(func(r chi.Router) {
				r.Use(httpmw.RequireAuth(tokenService))

				identityHandler.ProtectedRoutes(r)
				reviewsHandler.ProtectedRoutes(r)
				ordersHandler.Routes(r)
				checkoutHandler.Routes(r)
				favoritesHandler.Routes(r)
				loyaltyHandler.Routes(r)
				questionsHandler.ProtectedRoutes(r)
				exchangesHandler.Routes(r)

				// Seller (and any admin role): catalog management + sales analytics.
				r.Group(func(r chi.Router) {
					r.Use(httpmw.RequireRole(roleSeller, roleSystemAdmin, roleGlobalAdmin))
					catalogHandler.ManageRoutes(r)
					analyticsHandler.SellerRoutes(r)
					questionsHandler.SellerRoutes(r)
				})

				// System Admin + Global Admin: platform-wide metrics and palette config.
				r.Group(func(r chi.Router) {
					r.Use(httpmw.RequireRole(roleSystemAdmin, roleGlobalAdmin))
					analyticsHandler.AdminRoutes(r)
					settingsHandler.AdminRoutes(r)
				})

				// Global Admin only: infra-level health report.
				r.Group(func(r chi.Router) {
					r.Use(httpmw.RequireRole(roleGlobalAdmin))
					analyticsHandler.GlobalRoutes(r)
				})
			})
		})
	})

	// Health snapshot recorder: samples DB latency + platform counts
	// immediately at boot (so charts have a first data point right away)
	// and every 5 minutes thereafter. Feeds the Global Admin health
	// history charts with real, structured telemetry instead of mocks.
	go func() {
		if err := analyticsRepo.RecordSnapshot(context.Background()); err != nil {
			logger.Warn("failed to record initial health snapshot", zap.Error(err))
		}
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := analyticsRepo.RecordSnapshot(context.Background()); err != nil {
				logger.Warn("failed to record health snapshot", zap.Error(err))
			}
		}
	}()

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
