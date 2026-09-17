package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DB struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Schema   string `yaml:"schema"`
}

type Server struct {
	Addr string `yaml:"addr"`
}

type Config struct {
	DB DB `yaml:"db"`

	Logger struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"logger"`

	Server Server `yaml:"server"`
}

func Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("failed to validate config file: %w", err)
	}

	return &cfg, nil

}

func (c *Config) validate() error {
	if c.DB.Host == "" {
		return fmt.Errorf("db.host is required")
	}
	if c.DB.Port == 0 {
		return fmt.Errorf("db.port is required")
	}
	if c.DB.User == "" {
		return fmt.Errorf("db.user is required")
	}
	if c.DB.DBName == "" {
		return fmt.Errorf("db.dbname is required")
	}
	if c.DB.Schema == "" {
		return fmt.Errorf("db.schema is required")
	}
	if c.Logger.Level == "" {
		c.Logger.Level = "info"
	}
	if c.Logger.Format == "" {
		c.Logger.Format = "json"
	}
	if c.Server.Addr == "" {
		c.Server.Addr = ":8090"
	}
	return nil
}
