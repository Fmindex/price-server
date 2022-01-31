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
	GetPrices(currencies []string) (map[string]float64, error)
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
	currencies = []string{"BTC", "ETH", "LUNA"}
)

type GetLatestPriceResponse struct {
	Prices map[string]string `json:"prices"`
}

// Business Logic
// GetLatestPrice is to get the latest price of currencies from each SDKs
func (p *pricingImpl) GetLatestPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var priceListByCurrency = map[string][]float64{}

	// get price and then add to the list
	for _, exchangeSDK := range p.exchangeSDKs {
		prices, err := exchangeSDK.GetPrices(currencies)
		if err != nil {
			errno.GenErrorResp(w, errno.InternalError, err.Error())
			return
		}
		for _, currency := range currencies {
			priceListByCurrency[currency] = append(priceListByCurrency[currency], prices[currency])
		}
	}

	// find median from each currencies
	// sort array
	for _, currency := range currencies {
		sort.Float64s(priceListByCurrency[currency])
	}
	// find median
	var medianPriceByCurrency = map[string]string{}
	for _, currency := range currencies {
		priceLen := len(priceListByCurrency[currency])
		var price float64
		if priceLen%2 == 0 {
			// even, average of the middles
			price = (priceListByCurrency[currency][priceLen/2-1] + priceListByCurrency[currency][priceLen/2]) / 2.0
		} else {
			// odd, use the middle one
			price = priceListByCurrency[currency][priceLen/2]
		}
		medianPriceByCurrency[currency] = fmt.Sprintf("%f", price)
	}

	json.NewEncoder(w).Encode(GetLatestPriceResponse{
		Prices: medianPriceByCurrency,
	})
}
