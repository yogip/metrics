package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type ServerConfig struct {
	Address  string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL"`
}

type StorageConfig struct {
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	StoreIntreval   int64  `env:"STORE_INTERVAL"`
	Restore         bool   `env:"RESTORE"`
}

type Config struct {
	Server  ServerConfig
	HashKey string `env:"KEY"`
	Storage StorageConfig
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	flag.StringVar(&cfg.Server.Address, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.Server.LogLevel, "l", "info", "Log levle: debug, info, warn, error, panic, fatal")
	flag.Int64Var(&cfg.Storage.StoreIntreval, "i", 300, "Dump DB to file with given interval. 0 - means to write all changes immediately")
	flag.StringVar(&cfg.Storage.FileStoragePath, "f", "/tmp/metrics-db.json", "Path to dump file")
	flag.BoolVar(&cfg.Storage.Restore, "r", true, "Restore DB dump from file")
	flag.StringVar(&cfg.Storage.DatabaseDSN, "d", "", "Database connection string")
	flag.StringVar(&cfg.HashKey, "k", "", "Hash key to check request signature")

	flag.Parse()

	if value, exists := os.LookupEnv("ADDRESS"); exists {
		cfg.Server.Address = value
	}
	if value, exists := os.LookupEnv("LOG_LEVEL"); exists {
		cfg.Server.LogLevel = value
	}
	if value, exists := os.LookupEnv("STORE_INTERVAL"); exists {
		interval, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("STORE_INTERVAL convertation error: %w", err)
		}
		cfg.Storage.StoreIntreval = interval
	}
	if value, exists := os.LookupEnv("FILE_STORAGE_PATH"); exists {
		cfg.Storage.FileStoragePath = value
	}
	if value, exists := os.LookupEnv("RESTORE"); exists {
		restore, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("RESTORE convertation error: %w", err)
		}
		cfg.Storage.Restore = restore
	}
	if value, exists := os.LookupEnv("DATABASE_DSN"); exists {
		cfg.Storage.DatabaseDSN = value
	}
	if value, exists := os.LookupEnv("KEY"); exists && value != "" {
		cfg.HashKey = value
	}

	return &cfg, nil
}
