package binance

// mapper for binance synbols
var coinMapper = map[string]string{
	"BTC":  "BTCUSD",
	"LUNA": "LUNAUSD",
	"USD":  "BTCUSD",
}

// Init
type BinanceSDKImpl struct{}

func NewBinanceSDK() *BinanceSDKImpl {
	return &BinanceSDKImpl{}
}

func (b *BinanceSDKImpl) GetPrices(coins []string) (map[string]float64, error) {
	return map[string]float64{
		"BTC":  1,
		"LUNA": 1,
		"USD":  1,
	}, nil
}
