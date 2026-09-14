package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	} ` yaml:"db"`

	Logger struct {
		Level  string `json:"level"`
		Format string `json:"format"`
	} ` yaml:"logger"`

	Server struct {
		Addr string `json:"addr"`
	} ` yaml:"server"`
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
