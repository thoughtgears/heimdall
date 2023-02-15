package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug      bool   `envconfig:"DEBUG" default:"false"`
	Port       int    `envconfig:"PORT" default:"8080"`
	Project    string `envconfig:"PROJECT_ID" required:"true"`
	Collection string `envconfig:"COLLECTION" default:"projects"`
}

func New() (*Config, error) {
	var config *Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, fmt.Errorf("could not load config : %v", err)
	}

	return config, nil
}
