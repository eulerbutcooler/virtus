package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/eulerbutcooler/virtus/internal/config"
	"github.com/eulerbutcooler/virtus/internal/repository/postgres"
	dbgen "github.com/eulerbutcooler/virtus/internal/repository/postgres/db"
	redisrepo "github.com/eulerbutcooler/virtus/internal/repository/redis"
	"github.com/eulerbutcooler/virtus/internal/service"
	"github.com/eulerbutcooler/virtus/internal/worker/tasks"
	"github.com/eulerbutcooler/virtus/pkg/logger"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// Repositories and services needed by tasks.
	queries := dbgen.New(pool)
	cache := redisrepo.NewCache(rdb)

	userRepo := postgres.NewUserRepository(queries)
	poolRepo := postgres.NewPoolRepository(queries)
	requestRepo := postgres.NewRequestRepository(queries)
	queueRepo := postgres.NewQueueRepository(queries)

	authSvc := service.NewAuthService(userRepo, cache, cfg.JWT)
	_ = authSvc // not used by tasks; kept for potential future use

	poolSvc := service.NewPoolService(poolRepo)
	queueSvc := service.NewQueueService(queueRepo)
	requestSvc := service.NewRequestService(requestRepo, queueSvc)

	// Tasks
	fundTask := tasks.NewFundQueueTask(poolSvc, queueSvc, requestSvc)
	reminderTask := tasks.NewImpactReminderTask(pool)
	watchdogTask := tasks.NewDeliveryWatchdogTask(pool)

	slog.Info("worker starting")

	var wg sync.WaitGroup

	// Fund queue: runs every 30 seconds.
	// This is deliberately frequent so newly-completed contributions are acted on quickly.
	wg.Add(1)
	go func() {
		defer wg.Done()
		runPeriodic(ctx, 30*time.Second, "fund_queue", fundTask.Run)
	}()

	// Impact reminder: runs once per day.
	wg.Add(1)
	go func() {
		defer wg.Done()
		runPeriodic(ctx, 24*time.Hour, "impact_reminder", reminderTask.Run)
	}()

	// Delivery watchdog: runs once per day.
	wg.Add(1)
	go func() {
		defer wg.Done()
		runPeriodic(ctx, 24*time.Hour, "delivery_watchdog", watchdogTask.Run)
	}()

	// Block until interrupt.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	slog.Info("shutdown signal received", "signal", sig.String())

	cancel()  // stop all task goroutines
	wg.Wait() // wait for in-progress runs to finish

	slog.Info("worker stopped cleanly")
	return nil
}

// runPeriodic runs fn immediately and then on every tick of interval.
// It logs errors but does not stop — a failing task is retried on the next tick.
// It exits when ctx is cancelled.
func runPeriodic(ctx context.Context, interval time.Duration, name string, fn func(context.Context) error) {
	run := func() {
		if err := fn(ctx); err != nil {
			slog.Error("task error", "task", name, "err", err)
		}
	}

	// Run once immediately so the worker is useful from the first second.
	run()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			run()
		}
	}
}
