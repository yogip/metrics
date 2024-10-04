package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
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
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

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
		log.Fatalf("failed to initialize public key: %s", err)
		return
	}

	var wg sync.WaitGroup
	agent.Run(ctx, &wg, cfg, pubKey)

	<-quit
	logger.Log.Info("Received Ctrl+C, stopping...")
	cancel()
	wg.Wait() // wait for all goroutines to finish
}
