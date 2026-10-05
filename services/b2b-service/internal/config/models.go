package config

import (
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

type Config struct {
	LogLevel  string    `env:"LOG_LEVEL"`
	DemoMode  bool      `env:"DEMO_MODE"`
	HTTP      HTTP      `envPrefix:"HTTP_"`
	Postgres  Postgres  `envPrefix:"POSTGRES_"`
	Worker    Worker    `envPrefix:"WORKER_"`
	AIService AIService `envPrefix:"AI_SERVICE_"`
	MIS       MIS       `envPrefix:"MIS_"`
	Clinic    Clinic    `envPrefix:"CLINIC_"`
}

type HTTP struct {
	Addr              string        `env:"ADDR"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT"`
}

type Postgres = postgres.Config

type Worker = jobs.Config

type AIService struct {
	Mock    bool          `env:"MOCK"`
	URL     string        `env:"URL"`
	Timeout time.Duration `env:"TIMEOUT"`
}

type MIS struct {
	URL     string        `env:"URL"`
	Timeout time.Duration `env:"TIMEOUT"`
}

type Clinic struct {
	Phone string `env:"PHONE"`
}
