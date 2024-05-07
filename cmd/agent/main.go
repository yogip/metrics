package main

import (
	"flag"
	"log"
	"time"

	"metrics/internal/agent"
	"metrics/internal/agent/config"
)

var flagRunAddr string

var flagReportInterval int64
var flagPollInterval int64

func init() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "server addres and port to send metrics")
	flag.Int64Var(&flagReportInterval, "r", 10, "sent metric to server every given interval")
	flag.Int64Var(&flagPollInterval, "p", 2, "gather metric every given interval")
}

func main() {
	flag.Parse()

	cfg := config.NewAgentConfig(
		flagRunAddr,
		time.Duration(flagReportInterval)*time.Second,
		time.Duration(flagPollInterval)*time.Second,
	)

	log.Println("Start agent for: ", cfg.ServerAddresPort)
	agent.Run(cfg)
}
