package config

import "github.com/qreaqtor/wallet/pkg/logger"

type (
	Config struct {
		Log Logger `yaml:"logger"`
	}

	Logger struct {
		Level  logger.Level `yaml:"level" env-default:"info"`
		Pretty bool         `yaml:"pretty" env-default:"false"`
	}
)
