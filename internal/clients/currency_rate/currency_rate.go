package currency_rate

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type CurrencyRateApiClient interface {
	GetCurrencyRates() (map[string]ApiCurrencyRate, error)
}

type ApiCurrencyValute struct {
	Valute map[string]ApiCurrencyRate `json:"Valute"`
}

type ApiCurrencyRate struct {
	Name string  `json:"CharCode"`
	Rate float64 `json:"Value"`
}

type currencyRateApiClient struct {
	apiUrl string
}

func NewCurrencyRateApiClient(apiUrl string) CurrencyRateApiClient {
	return &currencyRateApiClient{
		apiUrl: apiUrl,
	}
}

func (c *currencyRateApiClient) GetCurrencyRates() (map[string]ApiCurrencyRate, error) {
	resp, err := http.Get(c.apiUrl)
	if err != nil {
		return nil, errors.Wrap(err, "api: getting currency_rate request error")
	}
	defer resp.Body.Close()

	var valute ApiCurrencyValute

	err = json.NewDecoder(resp.Body).Decode(&valute)
	if err != nil {
		return nil, errors.Wrap(err, "api: decoding response error")
	}

	return valute.Valute, nil
}
