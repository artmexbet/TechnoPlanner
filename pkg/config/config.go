package config

import "github.com/ilyakaznacheev/cleanenv"

func MustParseConfig[T any](configPath string) T {
	var res T
	err := cleanenv.ReadConfig(configPath, &res)
	if err != nil {
		panic(err)
	}
	return res
}
