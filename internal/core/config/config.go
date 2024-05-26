package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type ServerConfig struct {
	Address  string `env:"ADDRESS" envDefault:"localhost:8080"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

type StorageConfig struct {
	StoreIntreval   int64  `env:"STORE_INTERVAL" envDefault:"300"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/metrics-db.json"`
	Restore         bool   `env:"RESTORE" envDefault:"true"`
}

type Config struct {
	Server  ServerConfig
	Storage StorageConfig
}

func NewConfig() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	flag.StringVar(
		&cfg.Server.Address,
		"a",
		cfg.Server.Address,
		"address and port to run server",
	)
	flag.StringVar(
		&cfg.Server.LogLevel,
		"l",
		cfg.Server.LogLevel,
		"Log levle: debug, info, warn, error, panic, fatal",
	)
	flag.Int64Var(
		&cfg.Storage.StoreIntreval,
		"i",
		cfg.Storage.StoreIntreval,
		"Dump DB to file with given interval. 0 - means to write all changes immediately",
	)
	flag.StringVar(
		&cfg.Storage.FileStoragePath,
		"f",
		cfg.Storage.FileStoragePath,
		"Path to dump file",
	)
	flag.BoolVar(
		&cfg.Storage.Restore,
		"r",
		cfg.Storage.Restore,
		"Restore DB dump from file",
	)
	flag.Parse()

	return &cfg, nil
}
