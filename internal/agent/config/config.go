package config

import (
	"strings"
	"time"
)

type AgentConfig struct {
	ServerAddresPort string
	ReportInterval   time.Duration
	PollInterval     time.Duration
}

func NewAgentConfig(serverAddresPort string, reportInterval time.Duration, pollInterval time.Duration) *AgentConfig {
	if !strings.HasPrefix(serverAddresPort, "http://") && !strings.HasPrefix(serverAddresPort, "https://") {
		serverAddresPort = "http://" + serverAddresPort
	}

	return &AgentConfig{
		ServerAddresPort: serverAddresPort,
		ReportInterval:   reportInterval,
		PollInterval:     pollInterval,
	}
}
