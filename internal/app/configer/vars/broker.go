package vars

import "time"

type StreamConfig struct {
	Name     string        `yaml:"name"`
	Subjects []string      `yaml:"subjects"`
	MaxAge   time.Duration `yaml:"max_age"`
}

type Streams struct {
	Something StreamConfig `yaml:"something"`
}

type BrokerConfig struct {
	URL     string  `yaml:"url" env:"URL"`
	Streams Streams `yaml:"streams" env:"STREAMS"`
}
