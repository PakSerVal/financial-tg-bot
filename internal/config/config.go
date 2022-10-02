package config

import (
	"flag"
)

const defaultToken = "5332081649:AAHXDWuBHPleSVXzo8w8j2L8NaqbOll7B34"

type Config struct {
	token string
}

func New() (*Config, error) {
	c := &Config{}

	flag.StringVar(&c.token, "token", defaultToken, "bot token")
	flag.Parse()

	return c, nil
}

func (c *Config) Token() string {
	return c.token
}
