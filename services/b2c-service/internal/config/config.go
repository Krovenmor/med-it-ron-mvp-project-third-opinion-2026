package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

func Load() (Config, error) {
	cfg, err := env.ParseAsWithOptions[Config](env.Options{RequiredIfNoDef: true})
	if err != nil {
		return Config{}, fmt.Errorf("parse env: %w", err)
	}
	return cfg, nil
}
