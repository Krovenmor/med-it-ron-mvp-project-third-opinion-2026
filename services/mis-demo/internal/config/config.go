package config

import (
	"errors"
	"fmt"

	"github.com/caarlos0/env/v11"
)

func Load() (Config, error) {
	cfg, err := env.ParseAsWithOptions[Config](env.Options{RequiredIfNoDef: true})
	if err != nil {
		return Config{}, fmt.Errorf("parse env: %w", err)
	}
	if cfg.Schedule.SlotsCount < 1 {
		return Config{}, errors.New("SCHEDULE_SLOTS_COUNT must be positive")
	}
	return cfg, nil
}
