package common

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// Configs

type Log struct {
	Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
}

type HTTP struct {
	Port string `env-required:"true" yaml:"port" env:"SERVER_PORT"`
}

type Database struct {
	Host        string `env-required:"true" yaml:"host" env:"DATABASE_HOST"`
	Port        string `env-required:"true" yaml:"port" env:"DATABASE_PORT"`
	User        string `env-required:"true" yaml:"user" env:"DATABASE_USER"`
	Password    string `env-required:"true" yaml:"password" env:"DATABASE_PASSWORD"`
	Database    string `env-required:"true" yaml:"database" env:"DATABASE_NAME"`
	MaxConn     int    `env-required:"true" yaml:"max_conn" env:"DATABASE_MAX_CONN"`
	MaxIdleConn int    `env-required:"true" yaml:"max_idle_conn" env:"DATABASE_MAX_IDLE_CONN"`
	MaxLifetime int    `env-required:"true" yaml:"max_lifetime_conn" env:"DATABASE_MAX_LIFETIME_CONN"`
}

// DefaultConfigs

type Environment string

const (
	ModeDevelopment Environment = "development"
	ModeProduction  Environment = "production"
)

type Consul struct {
	Host string `env-required:"true" yaml:"host"`
	Port string `env-required:"true" yaml:"port"`
}

type Registry struct {
	Host string `env-required:"true" yaml:"host"`
}

type Config struct {
	HTTP     `yaml:"http"`
	Log      `yaml:"logger"`
	Database `yaml:"mysql"`
	Env      Environment `env-required:"true" yaml:"environment" env:"ENVIRONMENT"`
	Consul   Consul      `yaml:"consul"`
	Registry Registry    `yaml:"registry"`
}

func NewConfigFromYAMLFile(yamlFile string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(yamlFile, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
