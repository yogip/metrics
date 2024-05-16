package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"metrics/internal/core/service"
	"metrics/internal/infra/api/rest"
	"metrics/internal/infra/store"
	"metrics/internal/logger"

	"go.uber.org/zap"
)

func init() {

}

func main() {
	var runAddress string
	var logLevel string

	// todo move to config
	flag.StringVar(&runAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&logLevel, "l", "info", "Log levle: debug, info, warn, error, panic, fatal")
	flag.Parse()

	if address, ok := os.LookupEnv("ADDRESS"); ok {
		runAddress = address
	}
	if level, ok := os.LookupEnv("LOG_LEVEL"); ok {
		logLevel = level
	}

	err := logger.Initialize(logLevel)
	if err != nil {
		log.Fatal(err)
	}

	if err := run(runAddress); err != nil {
		logger.Log.Fatal("Running server Error", zap.String("error", err.Error()))
	}
}

func run(runAddress string) error {
	store, err := store.NewStore()
	if err != nil {
		return fmt.Errorf("failed to initialize a store: %w", err)
	}
	service := service.NewMetricService(store)
	logger.Log.Info("Service initialized")
	api := rest.NewAPI(service)

	logger.Log.Info("Start Server at:", zap.String("addres", runAddress))
	if err := api.Run(runAddress); err != nil {
		return fmt.Errorf("server has failed: %w", err)
	}
	return nil
}
