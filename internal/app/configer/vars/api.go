package vars

type APIConfig struct {
	Port int    `yaml:"port" env:"PORT"`
	Host string `yaml:"host" env:"HOST"`
}
