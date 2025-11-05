package config

import (
	"github.com/qreaqtor/wallet/pkg/logger"
)

type (
	Config struct {
		Port       int64 `env:"PORT" env-default:"8080"`
		LimiterRPS int64 `env:"LIMITER_RPS" env-default:"1000"`
		Log        Logger
		Database
	}

	Logger struct {
		Level  logger.Level `env:"LOGGER_LEVEL" env-default:"info"`
		Pretty bool         `env:"LOGGER_PRETTY" env-default:"false"`
	}

	Database struct {
		User         string `env:"DATABASE_USER" env-required:"true"`
		Password     string `env:"DATABASE_PASSWORD" env-required:"true"`
		DatabaseName string `env:"DATABASE_NAME" env-required:"true"`
		Address      string `env:"DATABASE_ADDRESS" env-required:"true"`
	}
)
