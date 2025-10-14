package vars

import "time"

type PostgresConfig struct {
	User string `yaml:"user" env:"USER"`
	Pass string `yaml:"password" env:"PASSWORD"`
	Host string `yaml:"host" env:"HOST"`
	Port string `yaml:"port" env:"PORT"`
	Db   string `yaml:"db" env:"DB"`

	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT"`
}
