package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"metrics/internal/core/config"
	"metrics/internal/core/service"
	"metrics/internal/infra/api/rest"
	"metrics/internal/infra/store"
	dbStorage "metrics/internal/infra/store/db"
	"metrics/internal/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = logger.Initialize(cfg.Server.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	if err := run(cfg); err != nil {
		logger.Log.Fatal("Running server Error", zap.String("error", err.Error()))
	}
}

func run(cfg *config.Config) error {
	ctx := context.Background()
	store, err := store.NewStore(&cfg.Storage)
	if err != nil {
		return fmt.Errorf("failed to initialize a store: %w", err)
	}

	dbStore, err := dbStorage.NewStore(&cfg.Storage)
	if err != nil {
		return fmt.Errorf("failed to initialize a store: %w", err)
	}

	metricService := service.NewMetricService(store)
	systemService := service.NewSystemService(dbStore)
	logger.Log.Info("Service initialized")
	api := rest.NewAPI(metricService, systemService)

	// https://github.com/gin-gonic/gin/blob/master/docs/doc.md#manually
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := api.Run(cfg.Server.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Info("Runing server error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := api.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	store.Close()
	dbStore.Close()
	logger.Log.Info("Server exiting")
	return nil
}
