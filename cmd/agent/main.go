package main

import (
	"flag"
	"log"
	"time"

	"metrics/internal/agent"
	"metrics/internal/agent/config"
)

var flagRunAddr string
var flagReportInterval time.Duration
var flagPollInterval time.Duration

func init() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "server addres and port to send metrics")
	flag.DurationVar(&flagReportInterval, "r", 10*time.Second, "sent metric to server every given interval")
	flag.DurationVar(&flagPollInterval, "p", 2*time.Second, "gather metric every given interval")
}

func main() {
	flag.Parse()

	cfg := config.NewAgentConfig(flagRunAddr, flagReportInterval, flagPollInterval)

	log.Println("Start agent for: ", cfg.ServerAddresPort)
	agent.Run(cfg)
}
