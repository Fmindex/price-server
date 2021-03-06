package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// mapper for binance synbols
var symbolMapper = map[string]string{
	"BTC":  "BTCUSDT",
	"LUNA": "LUNAUSDT",
	"ETH":  "ETHUSDT",
}

// Init
type BinanceSDKImpl struct{}

func NewBinanceSDK() *BinanceSDKImpl {
	return &BinanceSDKImpl{}
}

const (
	binanceHost  = "https://api.binance.com/api/v3/%s"
	getPricePath = "ticker/price?symbol=%s"
)

// A Response struct to map the Entire Response
type Response struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (b *BinanceSDKImpl) GetPrices(symbols []string) (map[string]float64, error) {
	var priceListBySymbol = map[string]float64{}
	for _, symbol := range symbols {

		// get binance symbol
		binanceSymbol, found := symbolMapper[symbol]
		if !found {
			// IMPROVEMENT: add logs
			continue
		}

		// get price
		path := fmt.Sprintf(getPricePath, binanceSymbol)
		response, err := http.Get(fmt.Sprintf(binanceHost, path))
		if err != nil {
			return nil, err
		}

		// parse response
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		var responseObject Response
		json.Unmarshal(responseData, &responseObject)

		// convert string to float
		f := responseObject.Price
		s, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return nil, err
		}

		// add price
		priceListBySymbol[symbol] = s
	}

	return priceListBySymbol, nil
}
