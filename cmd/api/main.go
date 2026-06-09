package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eulerbutcooler/virtus/internal/config"
	"github.com/eulerbutcooler/virtus/internal/handler"
	"github.com/eulerbutcooler/virtus/internal/repository/postgres"
	dbgen "github.com/eulerbutcooler/virtus/internal/repository/postgres/db"
	redisrepo "github.com/eulerbutcooler/virtus/internal/repository/redis"
	"github.com/eulerbutcooler/virtus/internal/service"
	"github.com/eulerbutcooler/virtus/pkg/logger"
	stripepkg "github.com/eulerbutcooler/virtus/pkg/stripe"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger.New(logger.Options{
		Level:  logger.LevelFromString(cfg.Log.Level),
		Format: cfg.Log.Format,
	})

	ctx := context.Background()

	// Postgres
	pool, err := postgres.New(ctx, postgres.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	defer pool.Close()

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("connect redis: %w", err)
	}
	defer rdb.Close()

	// Repositories
	queries := dbgen.New(pool)
	cache := redisrepo.NewCache(rdb)

	userRepo := postgres.NewUserRepository(queries)
	poolRepo := postgres.NewPoolRepository(queries)
	requestRepo := postgres.NewRequestRepository(queries)
	queueRepo := postgres.NewQueueRepository(queries)
	contributionRepo := postgres.NewContributionRepository(queries)
	fulfillmentRepo := postgres.NewFulfillmentRepository(queries)
	deliveryRepo := postgres.NewDeliveryRepository(queries)
	impactRepo := postgres.NewImpactRepository(queries)
	institutionRepo := postgres.NewInstitutionRepository(queries)

	// Stripe
	stripeProvider := stripepkg.New(cfg.Stripe.SecretKey, cfg.Stripe.WebhookSecret)

	// Services
	authSvc := service.NewAuthService(userRepo, cache, cfg.JWT)
	userSvc:=service.NewUserService(userRepo)
	poolSvc := service.NewPoolService(poolRepo)
	queueSvc := service.NewQueueService(queueRepo)
	requestSvc := service.NewRequestService(requestRepo, queueSvc)
	contributionSvc := service.NewContributionService(contributionRepo, poolSvc, pool, queries, stripeProvider)
	fulfillmentSvc := service.NewFulfillmentService(fulfillmentRepo, requestRepo, poolSvc)
	deliverySvc := service.NewDeliveryService(deliveryRepo, fulfillmentSvc)
	impactSvc := service.NewImpactService(impactRepo, deliveryRepo)
	institutionSvc := service.NewInstitutionService(institutionRepo, poolSvc)

	// HTTP server
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: handler.NewRouter(handler.Services{
			Auth:         authSvc,
			User: userSvc,
			Pool:         poolSvc,
			Request:      requestSvc,
			Queue:        queueSvc,
			Contribution: contributionSvc,
			Fulfillment:  fulfillmentSvc,
			Delivery:     deliverySvc,
			Impact:       impactSvc,
			Institution:  institutionSvc,
			Stripe:       stripeProvider,
		}),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("server starting", "addr", srv.Addr, "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case sig := <-stop:
		slog.Info("shutdown signal received", "signal", sig.String())
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	slog.Info("server stopped cleanly")
	return nil
}
