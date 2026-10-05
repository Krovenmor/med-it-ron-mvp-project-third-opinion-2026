package config

import "time"

type Config struct {
	LogLevel string `env:"LOG_LEVEL"`
	HTTP     HTTP   `envPrefix:"HTTP_"`
	B2B      B2B    `envPrefix:"B2B_"`
	MIS      MIS    `envPrefix:"MIS_"`
	Clinic   Clinic `envPrefix:"CLINIC_"`
}

type HTTP struct {
	Addr              string        `env:"ADDR"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT"`
}

type B2B struct {
	URL     string        `env:"URL"`
	Timeout time.Duration `env:"TIMEOUT"`
}

type MIS struct {
	URL          string        `env:"URL"`
	Timeout      time.Duration `env:"TIMEOUT"`
	SourceSystem string        `env:"SOURCE_SYSTEM"`
}

type Clinic struct {
	Name    string `env:"NAME"`
	Phone   string `env:"PHONE"`
	Address string `env:"ADDRESS"`
}
