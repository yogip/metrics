package main

import (
	"fmt"
	"log"

	"metrics/internal/core/config"
	"metrics/internal/core/service"
	"metrics/internal/infra/api/rest"
	"metrics/internal/infra/store"
	"metrics/internal/logger"

	"go.uber.org/zap"
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
	store, err := store.NewStore(&cfg.Storage)
	if err != nil {
		return fmt.Errorf("failed to initialize a store: %w", err)
	}

	service := service.NewMetricService(store)
	logger.Log.Info("Service initialized")
	api := rest.NewAPI(service)

	logger.Log.Info("Start Server at:", zap.String("addres", cfg.Server.Address))
	err = api.Run(cfg.Server.Address)
	if err != nil {
		log.Println(err)
	}
	store.Close()
	logger.Log.Info("Server exiting")
	return nil
}
