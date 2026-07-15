package core_redis

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr     string        `envconfig:"ADDR" default:"localhost:6379"`
	Password string        `envconfig:"PASSWORD"`
	DB       int           `envconfig:"DB" default:"0"`
	TTL      time.Duration `envconfig:"TTL" default:"5m"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("REDIS", &config); err != nil {
		return Config{}, fmt.Errorf("could not process env vars for redis config: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get redis config: %w", err)
		panic(err)
	}
	return config
}
