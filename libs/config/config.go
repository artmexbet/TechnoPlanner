package config

import (
	"fmt"

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
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"password"`
	DBName   string `yaml:"dbname" env:"POSTGRES_DB" env-default:"requests_db"`
	SSLMode  string `yaml:"sslmode" env:"POSTGRES_SSLMODE" env-default:"disable"`
}

func (pg Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pg.Host, pg.Port, pg.User, pg.Password, pg.DBName, pg.SSLMode)
}

type Trace struct {
	Endpoint string `yaml:"endpoint" env:"TRACE_ENDPOINT"`
	Insecure bool   `yaml:"insecure" env:"TRACE_INSECURE" env-default:"true"`
}
