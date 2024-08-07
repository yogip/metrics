package config

import (
	"flag"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
)

type AgentConfig struct {
	ServerAddresPort string `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval   int64  `env:"REPORT_INTERVAL" envDefault:"10"`
	PollInterval     int64  `env:"POLL_INTERVAL" envDefault:"2"`
	LogLevel         string `env:"LOG_LEVEL" envDefault:"info"`
	HashKey          string `env:"KEY"`
	RateLimit        int    `env:"RATE_LIMIT" envDefault:"3"`
}

func NewAgentConfig() (*AgentConfig, error) {
	var flagRunAddr string
	var flagReportInterval int64
	var flagPollInterval int64
	var flagLogLevel string
	var flagHashKey string
	var flagRateLimit int

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "server addres and port to send metrics")
	flag.Int64Var(&flagReportInterval, "r", 10, "sent metric to server every given interval")
	flag.Int64Var(&flagPollInterval, "p", 2, "gather metric every given interval")
	flag.StringVar(&flagLogLevel, "v", "info", "Log levle: debug, info, warn, error, panic, fatal")
	flag.StringVar(&flagHashKey, "k", "", "Hash key to sign requests")
	flag.IntVar(&flagRateLimit, "l", 3, "Amount of parallel requests to server")
	flag.Parse()

	cfg := AgentConfig{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
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
	if _, ok := os.LookupEnv("LOG_LEVEL"); !ok && flagLogLevel != "" {
		cfg.LogLevel = flagLogLevel
	}

	if !strings.HasPrefix(cfg.ServerAddresPort, "http://") && !strings.HasPrefix(cfg.ServerAddresPort, "https://") {
		cfg.ServerAddresPort = "http://" + cfg.ServerAddresPort
	}
	if _, ok := os.LookupEnv("KEY"); !ok && flagHashKey != "" {
		cfg.HashKey = flagHashKey
	}
	if _, ok := os.LookupEnv("RATE_LIMIT"); !ok && flagRateLimit > 0 {
		cfg.RateLimit = flagRateLimit
	}

	return &cfg, nil
}
