package main

import (
	"log"

	"metrics/internal/agent"
	"metrics/internal/agent/config"
)

func main() {
	cfg := config.NewAgentConfig()

	log.Println("Start agent for: ", cfg.ServerAddresPort)
	agent.Run(cfg)
}
