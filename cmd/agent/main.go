package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"metrics/internal/agent"
	"metrics/internal/agent/config"
	"metrics/internal/core/service"
	"metrics/internal/logger"

	"go.uber.org/zap"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	cfg, err := config.NewAgentConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		log.Fatal(err)
		return
	}

	logger.Log.Info("Start agent", zap.String("server", cfg.ServerAddresPort))
	logger.Log.Info(fmt.Sprintf("Build version: %s", buildVersion))
	logger.Log.Info(fmt.Sprintf("Build date: %s", buildDate))
	logger.Log.Info(fmt.Sprintf("Build commit: %s", buildCommit))

	pubKey, err := service.NewPublicKey(cfg.CryptoKey)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize public key: %s", err))
		return
	}

	agent.Run(ctx, cfg, pubKey)

	<-quit
	logger.Log.Info("Received Ctrl+C, stopping...")
	cancel()
}
