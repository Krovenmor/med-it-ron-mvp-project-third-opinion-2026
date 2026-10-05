package config

import "time"

type Config struct {
	LogLevel string   `env:"LOG_LEVEL"`
	HTTP     HTTP     `envPrefix:"HTTP_"`
	B2B      B2B      `envPrefix:"B2B_"`
	Schedule Schedule `envPrefix:"SCHEDULE_"`
}

type HTTP struct {
	Addr              string        `env:"ADDR"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT"`
}

type B2B struct {
	URL     string        `env:"URL"`
	Timeout time.Duration `env:"TIMEOUT"`
}

type Schedule struct {
	SlotsCount int `env:"SLOTS_COUNT"`
}
