package config

import (
	"github.com/GiDiS/jb-chat/pkg/auth"
	"github.com/GiDiS/jb-chat/pkg/logger"
	"github.com/caarlos0/env"
	"golang.org/x/oauth2"
	"strings"
)

const StoreMemory = "memory"
const StorePostgres = "postgres"

type Config struct {
	PublicPort int    `env:"PORT"`
	DiagPort   int    `env:"DIAG_PORT"`
	Store      string `env:"STORE"` // Choose store type
	Seed       bool   `env:"SEED"`  // Seed with GoT dataset
	Metrics    bool   `env:"METRICS_ENABLED"`
	Pprof      bool   `env:"PPROF_ENABLED"`

	GoogleAuth oauth2.Config
}

func MustBuild(log logger.Logger) Config {
	cfg := Config{
		PublicPort: 8888,
		DiagPort:   8889,
		Store:      StoreMemory,
		Seed:       false,
		Metrics:    true,
		Pprof:      false,
	}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Config build failed: %v", err)
	}

	cfg.GoogleAuth = auth.GetConfig()

	return cfg
}

func (c Config) GetStore() string {
	store := strings.ToLower(c.Store)
	if store == StorePostgres {
		return StorePostgres
	}
	return StoreMemory
}
