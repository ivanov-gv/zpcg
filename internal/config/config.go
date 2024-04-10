package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

const (
	EnvironmentProdValue    = "PROD"
	EnvironmentPreProdValue = "PREPROD"
)

type Config struct {
	TelegramApiToken string `env:"TELEGRAM_APITOKEN,required"`
	Port             string `env:"PORT,required"`
	Environment      string `env:"ENVIRONMENT" envDefault:"PROD"`
}

func Load() (Config, error) {
	var _config = &Config{}
	err := env.Parse(_config)
	if err != nil {
		return Config{}, fmt.Errorf("env.Parse: %w", err)
	}
	// return
	return *_config, nil
}
