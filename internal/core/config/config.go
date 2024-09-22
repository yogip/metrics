package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

type ServerConfig struct {
	Address  string
	LogLevel string
}

type StorageConfig struct {
	FileStoragePath string
	DatabaseDSN     string
	StoreIntreval   int64
	Restore         bool
}

type Config struct {
	Server    ServerConfig
	HashKey   string
	CryptoKey string
	Storage   StorageConfig
}

type JsonConfig struct {
	Address         *string `json:"address,omitempty"`
	LogLevel        *string `json:"log_level,omitempty"`
	HashKey         *string `json:"key,omitempty"`
	CryptoKey       *string `json:"crypto_key,omitempty"`
	FileStoragePath *string `json:"file_storage_path,omitempty"`
	DatabaseDSN     *string `json:"database_dsn,omitempty"`
	StoreIntreval   *int64  `json:"store_interval,omitempty"`
	Restore         *bool   `json:"restore,omitempty"`
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

func NewConfig() (*Config, error) {
	// Create config with default values
	cfg := Config{
		Server: ServerConfig{
			Address:  "localhost:8080",
			LogLevel: "info",
		},
		Storage: StorageConfig{
			FileStoragePath: "/tmp/metrics-db.json",
			DatabaseDSN:     "",
			StoreIntreval:   300,
			Restore:         true,
		},
		HashKey:   "",
		CryptoKey: "",
	}

	// Read commant args to serparate variables
	var jsonCfgPath, jsonCfgPathFull string
	var serverAddress, serverLogLevel string
	var storageStoreIntreval int64
	var storageFileStoragePath, storageDatabaseDSN string
	var storageRestore bool
	var hashKey, cryptoKey string

	flag.StringVar(&serverAddress, "a", "", "address and port to run server")
	flag.StringVar(&serverLogLevel, "l", "", "Log levle: debug, info, warn, error, panic, fatal")
	flag.Int64Var(&storageStoreIntreval, "i", 0, "Dump DB to file with given interval. 0 - means to write all changes immediately")
	flag.StringVar(&storageFileStoragePath, "f", "", "Path to dump file")
	flag.BoolVar(&storageRestore, "r", false, "Restore DB dump from file")
	flag.StringVar(&storageDatabaseDSN, "d", "", "Database connection string")
	flag.StringVar(&hashKey, "k", "", "Hash key to check request signature")
	flag.StringVar(&cryptoKey, "crypto-key", "", "Path to private key")
	flag.StringVar(&jsonCfgPath, "с", "", "json configuration file")
	flag.StringVar(&jsonCfgPathFull, "config", "", "json configuration file")

	flag.Parse()

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
	if value, exists := os.LookupEnv("ADDRESS"); exists {
		cfg.Server.Address = value
	} else if serverAddress != "" {
		cfg.Server.Address = serverAddress
	} else if jsonCfg != nil && jsonCfg.Address != nil {
		cfg.Server.Address = *jsonCfg.Address
	}

	// LOG_LEVEL
	if value, exists := os.LookupEnv("LOG_LEVEL"); exists {
		cfg.Server.LogLevel = value
	} else if serverLogLevel != "" {
		cfg.Server.LogLevel = serverLogLevel
	} else if jsonCfg != nil && jsonCfg.Address != nil {
		cfg.Server.LogLevel = *jsonCfg.LogLevel
	}

	// STORE_INTERVAL
	if value, exists := os.LookupEnv("STORE_INTERVAL"); exists {
		interval, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("STORE_INTERVAL convertation error: %w", err)
		}
		cfg.Storage.StoreIntreval = interval
	} else if storageStoreIntreval != 0 {
		cfg.Storage.StoreIntreval = storageStoreIntreval
	} else if jsonCfg != nil && jsonCfg.StoreIntreval != nil {
		cfg.Storage.StoreIntreval = *jsonCfg.StoreIntreval
	}

	// FILE_STORAGE_PATH
	if value, exists := os.LookupEnv("FILE_STORAGE_PATH"); exists {
		cfg.Storage.FileStoragePath = value
	} else if storageFileStoragePath != "" {
		cfg.Storage.FileStoragePath = storageFileStoragePath
	} else if jsonCfg != nil && jsonCfg.FileStoragePath != nil {
		cfg.Storage.FileStoragePath = *jsonCfg.FileStoragePath
	}

	// RESTORE
	if value, exists := os.LookupEnv("RESTORE"); exists {
		restore, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("RESTORE convertation error: %w", err)
		}
		cfg.Storage.Restore = restore
	} else if storageRestore {
		cfg.Storage.Restore = storageRestore
	} else if jsonCfg != nil && jsonCfg.FileStoragePath != nil {
		cfg.Storage.Restore = *jsonCfg.Restore
	}

	// DATABASE_DSN
	if value, exists := os.LookupEnv("DATABASE_DSN"); exists {
		cfg.Storage.DatabaseDSN = value
	} else if storageDatabaseDSN != "" {
		cfg.Storage.DatabaseDSN = storageDatabaseDSN
	} else if jsonCfg != nil && jsonCfg.FileStoragePath != nil {
		cfg.Storage.DatabaseDSN = *jsonCfg.DatabaseDSN
	}

	// KEY
	if value, exists := os.LookupEnv("KEY"); exists && value != "" {
		cfg.HashKey = value
	} else if hashKey != "" {
		cfg.HashKey = hashKey
	} else if jsonCfg != nil && jsonCfg.HashKey != nil {
		cfg.HashKey = *jsonCfg.HashKey
	}

	// CRYPTO_KEY
	if value, exists := os.LookupEnv("CRYPTO_KEY"); exists && value != "" {
		cfg.CryptoKey = value
	} else if cryptoKey != "" {
		cfg.CryptoKey = cryptoKey
	} else if jsonCfg != nil && jsonCfg.CryptoKey != nil {
		cfg.CryptoKey = *jsonCfg.CryptoKey
	}

	return &cfg, nil
}

func (cfg *Config) ReadCryptoKey() (*rsa.PrivateKey, error) {
	if cfg.CryptoKey == "" {
		return nil, nil
	}
	data, err := os.ReadFile(cfg.CryptoKey)
	if err != nil {
		return nil, fmt.Errorf("reading CRYPTO_KEY error: %w", err)
	}

	block, _ := pem.Decode(data)
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing CRYPTO_KEY error: %w", err)
	}

	return privateKey.(*rsa.PrivateKey), err
}
