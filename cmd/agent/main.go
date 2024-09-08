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
	}

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		log.Fatal(err)
	}

	logger.Log.Info("Start agent", zap.String("server", cfg.ServerAddresPort))
	logger.Log.Info(fmt.Sprintf("Build version: %s", buildVersion))
	logger.Log.Info(fmt.Sprintf("Build date: %s", buildDate))
	logger.Log.Info(fmt.Sprintf("Build commit: %s", buildCommit))
	agent.Run(ctx, cfg)

	<-quit
	logger.Log.Info("Received Ctrl+C, stopping...")
	cancel()
}
