package config

import (
	"github.com/caarlos0/env"
	"jb_chat/pkg/logger"
)

type Config struct {
	PublicPort int  `env:"PORT"`
	DiagPort   int  `env:"DIAG_PORT"`
	Seed       bool `env:"SEED"` // Seed this GoT dataset
}

func MustBuild(log logger.Logger) Config {
	cfg := Config{
		PublicPort: 8888,
		DiagPort:   8889,
		Seed:       false,
	}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Config build failed: %v", err)
	}

	return cfg
}
