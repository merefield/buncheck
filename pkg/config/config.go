package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Load loads the configuration file from the given path
func Load(path string) (*Configuration, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}

	var cfg Configuration
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct, %v", err)
	}

	return &cfg, nil
}

// Configuration holds application configuration data
type Configuration struct {
	DB DatabaseEnv `yaml:"database,omitempty"`
}

// DatabaseEnv holds dev and test database data
type DatabaseEnv struct {
	Dev  Database `yaml:"dev,omitempty"`
	Test Database `yaml:"test,omitempty"`
}

// Database holds data necessery for database configuration
type Database struct {
	PSN            string `yaml:"psn,omitempty"`
	LogQueries     bool   `yaml:"log_queries,omitempty"`
	TimeoutSeconds int    `yaml:"timeout_seconds,omitempty"`
}
