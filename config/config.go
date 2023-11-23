package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App  `yaml:"app"`
	HTTP `yaml:"http"`
	PG   `yaml:"postgres"`
	Nats `yaml:"nats"`
}

type App struct {
	Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
	Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
}

type HTTP struct {
	Host string `env-required:"true" yaml:"host" env:"HTTP_HOST"`
	Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
}

type PG struct {
	Host     string `env-required:"true" yaml:"host" env:"PG_HOST"`
	Port     string `env-required:"true" yaml:"port" env:"PG_PORT"`
	User     string `env-required:"true" yaml:"user" env:"PG_USER"`
	Password string `env-required:"true" yaml:"password" env:"PG_PASSWORD"`
	PgName   string `env-required:"true" yaml:"name" env:"PG_NAME"`
	PgDriver string `env-required:"true" yaml:"pg_driver" env:"PG_PG_DRIVER"`
}

type Nats struct {
	Host    string `env-required:"true" yaml:"host" env:"NATS_HOST"`
	Port    string `env-required:"true" yaml:"port" env:"NATS_PORT"`
	Cluster string `env-required:"true" yaml:"cluster" env:"NATS_CLUSTER"`
	Client  string `env-required:"true" yaml:"client" env:"NATS_CLIENT"`
	Topic   string `env-required:"true" yaml:"topic" env:"NATS_TOPIC"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("config.yml", cfg) //  !!!!  в гите тут исправить "config.yml"
	if err != nil {
		return nil, fmt.Errorf("Config error: %v", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("Environment reading error: %v", err)

	}

	return cfg, nil
}
