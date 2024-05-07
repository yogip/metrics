package config

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
)

type AgentConfig struct {
	ServerAddresPort string `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval   int64  `env:"REPORT_INTERVAL" envDefault:"10"`
	PollInterval     int64  `env:"POLL_INTERVAL" envDefault:"2"`
}

func NewAgentConfig() *AgentConfig {
	var flagRunAddr string
	var flagReportInterval int64
	var flagPollInterval int64

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "server addres and port to send metrics")
	flag.Int64Var(&flagReportInterval, "r", 10, "sent metric to server every given interval")
	flag.Int64Var(&flagPollInterval, "p", 2, "gather metric every given interval")
	flag.Parse()

	cfg := AgentConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Panicf("config parsing error: %s", err)
	}

	if _, ok := os.LookupEnv("ADDRESS"); !ok && flagRunAddr != "" {
		cfg.ServerAddresPort = flagRunAddr
	}
	if _, ok := os.LookupEnv("REPORT_INTERVAL"); !ok && flagReportInterval > 0 {
		cfg.ReportInterval = flagReportInterval
	}
	if _, ok := os.LookupEnv("POLL_INTERVAL"); !ok && flagPollInterval > 0 {
		cfg.PollInterval = flagPollInterval
	}

	if !strings.HasPrefix(cfg.ServerAddresPort, "http://") && !strings.HasPrefix(cfg.ServerAddresPort, "https://") {
		cfg.ServerAddresPort = "http://" + cfg.ServerAddresPort
	}

	return &cfg
}
