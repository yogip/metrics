package config

import (
	"time"
)

type AgentConfig struct {
	ServerAddresPort string
	ReportInterval   time.Duration
	PollInterval     time.Duration
}

func NewAgentConfig(serverAddresPort string, reportInterval time.Duration, pollInterval time.Duration) *AgentConfig {
	return &AgentConfig{
		ServerAddresPort: serverAddresPort,
		ReportInterval:   reportInterval,
		PollInterval:     pollInterval,
	}
}
