package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug                  bool   `envconfig:"LOCAL_DEBUG"`
	ProjectID              string `envconfig:"PROJECT_ID" required:"true"`
	Port                   string `envconfig:"PORT" required:"true"`
}

func BbwyConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
