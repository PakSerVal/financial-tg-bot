package config

import (
	"flag"
)

const (
	defaultBotToken            = "5332081649:AAEUJ8cGPsEaoAAFxwhqplbsDup2Rw2ez2s"
	defaultCurrencyRatesApiUrl = "https://www.cbr-xml-daily.ru/daily_json.js"
)

type Config struct {
	token               string
	currencyRatesApiUrl string
}

func New() (*Config, error) {
	c := &Config{}

	flag.StringVar(&c.token, "token", defaultBotToken, "bot token")
	flag.StringVar(&c.currencyRatesApiUrl, "currencyRatesApiUrl", defaultCurrencyRatesApiUrl, "currency_rate exchange url")
	flag.Parse()

	return c, nil
}

func (c *Config) Token() string {
	return c.token
}

func (c *Config) CurrencyRatesApiUrl() string {
	return c.currencyRatesApiUrl
}
