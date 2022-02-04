package coingecko

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// mapper for coingecko synbols
var symbolMapper = map[string]string{
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
	getPricePath  = "simple/price?ids=%s&vs_symbols=%s"
	vsSymbol      = "usd"
)

func (b *CoingeckoSDKImpl) GetPrices(symbols []string) (map[string]float64, error) {

	var priceListBySymbol = map[string]float64{}

	// get joined of mapped symbols
	var mappedSymbols []string
	for _, symbol := range symbols {
		// get coingecko symbol
		coingeckoSymbol, found := symbolMapper[symbol]
		if !found {
			// IMPROVEMENT: add logs
			continue
		}

		mappedSymbols = append(mappedSymbols, coingeckoSymbol)
	}

	// get price
	path := fmt.Sprintf(getPricePath, strings.Join(mappedSymbols, ","), vsSymbol)
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
	for _, symbol := range symbols {
		// get coingecko symbol
		coingeckoSymbol, found := symbolMapper[symbol]
		if !found {
			// IMPROVEMENT: add logs
			continue
		}

		// extract price from response
		f, found := responseObject[coingeckoSymbol][vsSymbol]
		if !found {
			// IMPROVEMENT: add logs
			continue
		}

		// add price
		priceListBySymbol[symbol] = f
	}

	return priceListBySymbol, nil
}
