package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/pkg/errors"
)

type AgentConfig struct {
	ServerAddresPort string `env:"ADDRESS" envDefault:"localhost:8080"`
	LogLevel         string `env:"LOG_LEVEL" envDefault:"info"`
	HashKey          string `env:"KEY"`
	CryptoKey        string `env:"CRYPTO_KEY"`
	ReportInterval   int64  `env:"REPORT_INTERVAL" envDefault:"10"`
	PollInterval     int64  `env:"POLL_INTERVAL" envDefault:"2"`
	RateLimit        int    `env:"RATE_LIMIT" envDefault:"3"`
}

type JsonConfig struct {
	ServerAddresPort *string `json:"address,omitempty"`
	LogLevel         *string `json:"log_level,omitempty"`
	HashKey          *string `json:"key,omitempty"`
	CryptoKey        *string `json:"crypto_key"`
	ReportInterval   *int64  `json:"report_interval"`
	PollInterval     *int64  `json:"poll_interval"`
	RateLimit        *int    `json:"rate_limit"`
}

func loadJsonConfig(path string) (cfg *JsonConfig, err error) {
	if path == "" {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "reading json config file error")
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, errors.Wrap(err, "unmarshal json config error")
	}
	return cfg, nil
}

func NewAgentConfig() (*AgentConfig, error) {
	// Read commant args to serparate variables
	var jsonCfgPath, jsonCfgPathFull string
	var flagRunAddr string
	var flagLogLevel string
	var flagHashKey string
	var flagCryptoKey string
	var flagReportInterval int64
	var flagPollInterval int64
	var flagRateLimit int

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "server addres and port to send metrics")
	flag.Int64Var(&flagReportInterval, "r", 10, "sent metric to server every given interval")
	flag.Int64Var(&flagPollInterval, "p", 2, "gather metric every given interval")
	flag.StringVar(&flagLogLevel, "v", "info", "Log levle: debug, info, warn, error, panic, fatal")
	flag.StringVar(&flagHashKey, "k", "", "Hash key to sign requests")
	flag.StringVar(&flagCryptoKey, "crypto-key", "", "Path to private key")
	flag.IntVar(&flagRateLimit, "l", 3, "Amount of parallel requests to server")
	flag.StringVar(&jsonCfgPath, "с", "", "json configuration file")
	flag.StringVar(&jsonCfgPathFull, "config", "", "json configuration file")
	flag.Parse()

	cfg := AgentConfig{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	// Build json config
	if jsonCfgPath == "" {
		jsonCfgPath = jsonCfgPathFull
	}
	if value, exists := os.LookupEnv("CONFIG"); exists {
		jsonCfgPath = value
	}
	jsonCfg, err := loadJsonConfig(jsonCfgPath)
	if err != nil {
		return nil, fmt.Errorf("json config loading error: %w", err)
	}

	// Read Env variables and override config values
	// ADDRESS
	if _, ok := os.LookupEnv("ADDRESS"); !ok && flagRunAddr != "" {
		cfg.ServerAddresPort = flagRunAddr
	} else if jsonCfg != nil && jsonCfg.ServerAddresPort != nil {
		cfg.ServerAddresPort = *jsonCfg.ServerAddresPort
	}
	if !strings.HasPrefix(cfg.ServerAddresPort, "http://") && !strings.HasPrefix(cfg.ServerAddresPort, "https://") {
		cfg.ServerAddresPort = "http://" + cfg.ServerAddresPort
	}

	// REPORT_INTERVAL
	if _, ok := os.LookupEnv("REPORT_INTERVAL"); !ok && flagReportInterval > 0 {
		cfg.ReportInterval = flagReportInterval
	} else if jsonCfg != nil && jsonCfg.ReportInterval != nil {
		cfg.ReportInterval = *jsonCfg.ReportInterval
	}

	// POLL_INTERVAL
	if _, ok := os.LookupEnv("POLL_INTERVAL"); !ok && flagPollInterval > 0 {
		cfg.PollInterval = flagPollInterval
	} else if jsonCfg != nil && jsonCfg.PollInterval != nil {
		cfg.PollInterval = *jsonCfg.PollInterval
	}

	// LOG_LEVEL
	if _, ok := os.LookupEnv("LOG_LEVEL"); !ok && flagLogLevel != "" {
		cfg.LogLevel = flagLogLevel
	} else if jsonCfg != nil && jsonCfg.LogLevel != nil {
		cfg.LogLevel = *jsonCfg.LogLevel
	}

	// KEY
	if _, ok := os.LookupEnv("KEY"); !ok && flagHashKey != "" {
		cfg.HashKey = flagHashKey
	} else if jsonCfg != nil && jsonCfg.HashKey != nil {
		cfg.HashKey = *jsonCfg.HashKey
	}

	// CRYPTO_KEY
	if _, ok := os.LookupEnv("CRYPTO_KEY"); !ok && flagCryptoKey != "" {
		cfg.CryptoKey = flagCryptoKey
	} else if jsonCfg != nil && jsonCfg.CryptoKey != nil {
		cfg.CryptoKey = *jsonCfg.CryptoKey
	}

	// RATE_LIMIT
	if _, ok := os.LookupEnv("RATE_LIMIT"); !ok && flagRateLimit > 0 {
		cfg.RateLimit = flagRateLimit
	} else if jsonCfg != nil && jsonCfg.RateLimit != nil {
		cfg.RateLimit = *jsonCfg.RateLimit
	}

	return &cfg, nil
}
