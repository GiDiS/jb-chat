package config

import (
	"github.com/GiDiS/jb-chat/pkg/logger"
	"github.com/caarlos0/env"
)

type Config struct {
	PublicPort int  `env:"PORT"`
	DiagPort   int  `env:"DIAG_PORT"`
	Seed       bool `env:"SEED"` // Seed this GoT dataset
	Metrics    bool `env:"METRICS_ENABLED"`
	Pprof      bool `env:"PPROF_ENABLED"`
}

func MustBuild(log logger.Logger) Config {
	cfg := Config{
		PublicPort: 8888,
		DiagPort:   8889,
		Seed:       false,
		Metrics:    true,
		Pprof:      false,
	}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Config build failed: %v", err)
	}

	return cfg
}
