package pricing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	errno "github.com/Fmindex/price-server/internal/pkg/error"
)

// Init
type ExchangeSDK interface {
	GetPrices(coins []string) (map[string]float64, error)
}

type pricingImpl struct {
	exchangeSDKs []ExchangeSDK
}

func NewPricing(e []ExchangeSDK) *pricingImpl {
	return &pricingImpl{exchangeSDKs: e}
}

// consts and declarations
var (
	// IMPROVEMENT: can maintain this list as a configuration separated from business logic
	// IMPROVEMENT: can make conis to be an enums
	// ASSUMPTION: this list is a price in USD
	coins = []string{"BTC", "ETH", "LUNA"}
)

type GetLatestPriceResponse struct {
	Prices map[string]string `json:"prices"`
}

// Business Logic
// GetLatestPrice is to get the latest price of coins from each SDKs
func (p *pricingImpl) GetLatestPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var priceListByCoin = map[string][]float64{}

	// get price and then add to the list
	for _, exchangeSDK := range p.exchangeSDKs {
		prices, err := exchangeSDK.GetPrices(coins)
		if err != nil {
			errno.GenErrorResp(w, errno.InternalError, err.Error())
			return
		}
		for _, coin := range coins {
			priceListByCoin[coin] = append(priceListByCoin[coin], prices[coin])
		}
	}

	// find median from each coins
	// sort array
	for _, coin := range coins {
		sort.Float64s(priceListByCoin[coin])
	}
	// find median
	var medianPriceByCoin = map[string]string{}
	for _, coin := range coins {
		priceLen := len(priceListByCoin[coin])
		var price float64
		if priceLen%2 == 0 {
			// even, average of the middles
			price = (priceListByCoin[coin][priceLen/2-1] + priceListByCoin[coin][priceLen/2]) / 2.0
		} else {
			// odd, use the middle one
			price = priceListByCoin[coin][priceLen/2]
		}
		medianPriceByCoin[coin] = fmt.Sprintf("%f", price)
	}

	json.NewEncoder(w).Encode(GetLatestPriceResponse{
		Prices: medianPriceByCoin,
	})
}
