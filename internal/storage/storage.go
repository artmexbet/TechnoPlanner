package storage

type Config struct {
	User string `yaml:"user" env:"USER"`
	Pass string `yaml:"password" env:"PASSWORD"`
	Host string `yaml:"host" env:"HOST"`
	Port string `yaml:"port" env:"PORT"`
	Db   string `yaml:"db" env:"DB"`
}

type Storage struct {
	cfg Config
}

func NewStorage(cfg Config) *Storage {
	return &Storage{
		cfg: cfg,
	}
}
