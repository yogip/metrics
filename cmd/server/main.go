package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"metrics/internal/core/config"
	"metrics/internal/core/service"
	"metrics/internal/infra/api/grpc"
	"metrics/internal/infra/api/rest"
	"metrics/internal/infra/store"
	"metrics/internal/logger"
	"metrics/migrations"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

type RunBackend interface {
	Run(string) error               // Run HTTP or GRPC backend
	Shutdown(context.Context) error // Shutdown backend
}

type runner struct {
	app   RunBackend
	cfg   *config.Config
	wg    *sync.WaitGroup
	store store.Store
}

func newRunner(ctx context.Context) (*runner, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = logger.Initialize(cfg.Server.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	logger.Log.Info(fmt.Sprintf("Build version: %s", buildVersion))
	logger.Log.Info(fmt.Sprintf("Build date: %s", buildDate))
	logger.Log.Info(fmt.Sprintf("Build commit: %s", buildCommit))

	if cfg.Storage.DatabaseDSN != "" {
		err := migrations.RunMigration(ctx, cfg)
		if err != nil {
			logger.Log.Fatal("Making migration error", zap.String("error", err.Error()))
		}
	}

	var wg sync.WaitGroup

	store, err := store.NewStore(
		ctx,
		&wg,
		&cfg.Storage,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a store: %w", err)
	}

	metricService := service.NewMetricService(store)
	systemService := service.NewSystemService(store)
	logger.Log.Info("Service initialized")

	privateKey, err := service.NewPrivateKey(cfg.CryptoKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize private key: %w", err)
	}

	r := runner{
		cfg:   cfg,
		wg:    &wg,
		store: store,
	}

	switch cfg.Server.BackendType {
	case config.HTTPBackType:
		r.app = rest.NewAPI(cfg, metricService, systemService, privateKey)
	case config.GRPCBackType:
		r.app = grpc.NewMetricsServer(cfg, metricService, systemService, privateKey)
	default:
		return nil, fmt.Errorf("unknown backend type: %s", cfg.Server.BackendType)
	}

	return &r, nil
}

// Run HTTP or GRPC backend.
func (r *runner) RunBackend() {
	err := r.app.Run(r.cfg.Server.Address)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Info("Runing server error", zap.Error(err))
	}
}

func (r *runner) Shutdown(ctx context.Context) error {
	if r.store != nil {
		r.store.Close()
	}
	if err := r.app.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// Run backend with gracefull shutdown.
func (r *runner) Run(ctx context.Context) error {
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go r.RunBackend()

	// Gracefully shutdown logic...(https://github.com/gin-gonic/gin/blob/master/docs/doc.md#manually)
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.Shutdown(ctx)
	if err != nil {
		return err
	}

	r.wg.Wait() // wait for all goroutines to finish
	return nil
}

func main() {
	ctx := context.Background()

	r, err := newRunner(ctx)
	if err != nil {
		logger.Log.Fatal("Running server Error", zap.String("error", err.Error()))
	}

	if err := r.Run(ctx); err != nil {
		logger.Log.Fatal("Running server Error", zap.String("error", err.Error()))
	}
	logger.Log.Info("Server exiting")
}
