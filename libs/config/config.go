package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

func MustParseConfig[T any](configPath string) T {
	var res T
	err := cleanenv.ReadConfig(configPath, &res)
	if err != nil {
		panic(err)
	}
	return res
}

type Postgres struct {
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	User     string `yaml:"user" env:"USER" env-default:"postgres"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"postgres"`
	DBName   string `yaml:"db_name" env:"DB"`
	SSLMode  string `yaml:"sslmode" env:"SSLMODE" env-default:"disable"`

	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"10s"`
}

func (pg Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pg.Host, pg.Port, pg.User, pg.Password, pg.DBName, pg.SSLMode)
}

type Trace struct {
	Endpoint string `yaml:"endpoint" env:"TRACE_ENDPOINT"`
	Insecure bool   `yaml:"insecure" env:"TRACE_INSECURE" env-default:"true"`
}

type NATSConfig struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"PORT" env-default:"4222"`
}

func (cfg NATSConfig) URL() string {
	return fmt.Sprintf("nats://%s:%s", cfg.Host, cfg.Port)
}
