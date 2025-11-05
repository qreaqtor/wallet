package di

import (
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/qreaqtor/wallet/internal/infrastucture/di/config"
)

func LoadConfig() (config.Config, error) {
	cur, err := os.Getwd()
	if err != nil {
		return config.Config{}, nil
	}

	cfgPath := "config.env"

	if strings.Contains(cur, "ac_tests") {
		cfgPath = strings.Split(cur, "ac_tests")[0] + cfgPath
	}

	cfg := new(config.Config)

	err = cleanenv.ReadConfig(cfgPath, cfg)
	if err != nil {
		return config.Config{}, err
	}

	return *cfg, nil
}
