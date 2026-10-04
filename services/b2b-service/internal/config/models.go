package config

import "time"

type Config struct {
	LogLevel  string    `env:"LOG_LEVEL"`
	DemoMode  bool      `env:"DEMO_MODE"`
	HTTP      HTTP      `envPrefix:"HTTP_"`
	Postgres  Postgres  `envPrefix:"POSTGRES_"`
	Worker    Worker    `envPrefix:"WORKER_"`
	AIService AIService `envPrefix:"AI_SERVICE_"`
	MIS       MIS       `envPrefix:"MIS_"`
}

type HTTP struct {
	Addr              string        `env:"ADDR"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT"`
}

type Postgres struct {
	DSN            string        `env:"DSN"`
	ConnectTimeout time.Duration `env:"CONNECT_TIMEOUT"`
}

type Worker struct {
	Count        int           `env:"COUNT"`
	PollInterval time.Duration `env:"POLL_INTERVAL"`
	Lease        time.Duration `env:"LEASE"`
	JobTimeout   time.Duration `env:"JOB_TIMEOUT"`
	MaxAttempts  int           `env:"MAX_ATTEMPTS"`
	BackoffBase  time.Duration `env:"BACKOFF_BASE"`
	BackoffMax   time.Duration `env:"BACKOFF_MAX"`
}

type AIService struct {
	Mock    bool          `env:"MOCK"`
	URL     string        `env:"URL"`
	Timeout time.Duration `env:"TIMEOUT"`
}

type MIS struct {
	URL     string        `env:"URL"`
	Timeout time.Duration `env:"TIMEOUT"`
}
