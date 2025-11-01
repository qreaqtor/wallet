package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	errNotExists = fmt.Errorf("file not exists")
)

const configFilename = "env.yaml"

func Load() (Config, error) {
	curDir, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}

	path := filepath.Join(curDir, configFilename)

	_, err = os.Stat(path)
	if err != nil {
		return Config{}, errNotExists
	}

	cfg := new(Config)

	err = cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return Config{}, err
	}

	return *cfg, nil
}
