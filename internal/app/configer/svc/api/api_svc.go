package api

import (
	"time"

	"technoBro/internal/app/configer/vars"
)

type Config struct {
	Postgres vars.PostgresConfig `yaml:"postgres" env:"POSTGRES"`
	Broker   vars.BrokerConfig   `yaml:"broker" env:"BROKER"`
	API      vars.APIConfig      `yaml:"api" env:"API"`
}

// Default возвращает конфигурацию по умолчанию для сервиса API.
func Default() Config {
	return Config{
		Postgres: vars.PostgresConfig{
			User:    "user",
			Pass:    "password",
			Host:    "localhost",
			Port:    "5432",
			Db:      "dbname",
			Timeout: 5 * time.Second,
		},
		Broker: vars.BrokerConfig{
			URL: "",
			Streams: vars.Streams{
				Something: vars.StreamConfig{
					Name:     "SMTHING",
					Subjects: []string{"smthing.created"},
					MaxAge:   24 * time.Hour, // 1d
				},
			},
		},
		API: vars.APIConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
	}
}
