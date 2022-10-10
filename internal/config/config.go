package config

import (
	"flag"
)

type Config struct {
	token string
}

func New() (*Config, error) {
	c := &Config{}

	flag.StringVar(&c.token, "token", "", "bot token")
	flag.Parse()

	return c, nil
}

func (c *Config) Token() string {
	return c.token
}
