package config

import (
	"time"

	"github.com/eddie023/byd/internal/build"
	"github.com/kelseyhightower/envconfig"
)

type APIConfig struct {
	Version string

	Web struct {
		ReadTimeout     time.Duration `default:"5s"`
		WriteTimeout    time.Duration `default:"10s"`
		IdleTimeout     time.Duration `default:"120s"`
		ShutdownTimeout time.Duration `default:"20s"`
		APIHost         string        `default:"0.0.0.0:8000" envconfig:"API_HOST"`
		DebugHost       string        `default:"0.0.0.0:4000"`
	}

	Db struct {
		ConnectionURI string `envconfig:"DB_CONNECTION_URI" required:"true"`
	}
}

// Returns parsed config
func New() (*APIConfig, error) {
	cfg := APIConfig{
		Version: build.Build,
	}

	err := envconfig.Process("api", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
