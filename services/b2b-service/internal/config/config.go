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
	if err := cfg.validate(); err != nil {
		return Config{}, fmt.Errorf("validate config: %w", err)
	}
	return cfg, nil
}

func (c Config) validate() error {
	switch {
	case c.Worker.Count < 1:
		return errors.New("WORKER_COUNT must be positive")
	case c.Worker.MaxAttempts < 1:
		return errors.New("WORKER_MAX_ATTEMPTS must be positive")
	case c.Worker.Lease <= c.Worker.JobTimeout:
		return errors.New("WORKER_LEASE must exceed WORKER_JOB_TIMEOUT, otherwise a running job can be claimed twice")
	case !c.AIService.Mock && c.AIService.URL == "":
		return errors.New("AI_SERVICE_URL is required when AI_SERVICE_MOCK is false")
	}
	return nil
}
