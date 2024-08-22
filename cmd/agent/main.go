package main

import (
	"log"

	"metrics/internal/agent"
	"metrics/internal/agent/config"
	"metrics/internal/logger"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.NewAgentConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		log.Fatal(err)
	}

	logger.Log.Info("Start agent", zap.String("server", cfg.ServerAddresPort))
	agent.Run(cfg)
}
