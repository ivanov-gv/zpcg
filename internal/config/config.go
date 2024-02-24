package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"

	"zpcg/resources"
)

type Config struct {
	TelegramApiToken     string `env:"TELEGRAM_APITOKEN,required"`
	Port                 string `env:"PORT,required"`
	TimetableGobFileName string `env:"TIMETABLE_GOB_FILENAME"`
}

func Load() (Config, error) {
	var _config *Config
	err := env.Parse(_config)
	if err != nil {
		return Config{}, fmt.Errorf("env.Parse: %w", err)
	}
	// set defaults
	if len(_config.TimetableGobFileName) == 0 {
		_config.TimetableGobFileName = resources.TimetableGobFileName
	}
	// return
	return *_config, nil
}
