package coingecko

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// mapper for coingecko synbols
var currencyMapper = map[string]string{
	"BTC":  "bitcoin",
	"LUNA": "terra-luna",
	"ETH":  "ethereum",
}

// Init
type CoingeckoSDKImpl struct{}

func NewCoingeckoSDK() *CoingeckoSDKImpl {
	return &CoingeckoSDKImpl{}
}

const (
	coingeckoHost = "https://api.coingecko.com/api/v3/%s"
	getPricePath  = "simple/price?ids=%s&vs_currencies=%s"
	vsCurrency    = "usd"
)

func (b *CoingeckoSDKImpl) GetPrices(currencies []string) (map[string]float64, error) {

	var priceListByCurrency = map[string]float64{}

	// get joined of mapped currencies
	var mappedCurrencies []string
	for _, currency := range currencies {
		// get coingecko symbol
		coingeckoSymbol, found := currencyMapper[currency]
		if !found {
			// IMPROVEMENT: add logs
			continue
		}

		mappedCurrencies = append(mappedCurrencies, coingeckoSymbol)
	}

	// get price
	path := fmt.Sprintf(getPricePath, strings.Join(mappedCurrencies, ","), vsCurrency)
	response, err := http.Get(fmt.Sprintf(coingeckoHost, path))
	if err != nil {
		return nil, err
	}

	// parse response
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var responseObject = map[string]map[string]float64{}
	json.Unmarshal(responseData, &responseObject)

	// build returned price object
	for _, currency := range currencies {
		// get coingecko symbol
		coingeckoSymbol, found := currencyMapper[currency]
		if !found {
			// IMPROVEMENT: add logs
			continue
		}

		// extract price from response
		f, found := responseObject[coingeckoSymbol][vsCurrency]
		if !found {
			// IMPROVEMENT: add logs
			continue
		}

		// add price
		priceListByCurrency[currency] = f
	}

	return priceListByCurrency, nil
}
