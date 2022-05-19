package main

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug                  bool   `envconfig:"LOCAL_DEBUG"`
	ProjectID              string `envconfig:"PROJECT_ID" required:"true"`
	Port                   string `envconfig:"PORT" required:"true"`
	DBConnString           string `envconfig:"POSTGRES_CONNSTR" required:"true"`
	BusinessMetricServer   string `envconfig:"BUSINESS_METRIC_SERVER" required:"true"`
	BusinessMetricToken    string `envconfig:"BUSINESS_METRIC_TOKEN" required:"true"`
	BusinessMetricOrg      string `envconfig:"BUSINESS_METRIC_ORG" required:"true"`
	BusinessMetricBucket   string `envconfig:"BUSINESS_METRIC_BUCKET" required:"true"`
}

func AdminConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
