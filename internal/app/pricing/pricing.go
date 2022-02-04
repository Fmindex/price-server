package pricing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
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
	wg         sync.WaitGroup
)

type GetLatestPriceResponse struct {
	Prices map[string]string `json:"prices"`
}

// Business Logic
// GetLatestPrice is to get the latest price of currencies from each SDKs
func (p *pricingImpl) GetLatestPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get price
	medianPriceByCurrency := p.getLatestPriceImpl()

	// api return
	json.NewEncoder(w).Encode(GetLatestPriceResponse{
		Prices: medianPriceByCurrency,
	})
}

func (p *pricingImpl) getLatestPriceImpl() map[string]string {

	var priceListByCurrency = map[string][]float64{}
	var priceByExchange = []map[string]float64{}

	// get price from each exchange sources asynchoronously
	wg.Add(len(p.exchangeSDKs))
	for _, exchangeSDK := range p.exchangeSDKs {
		var pricefromCurrentExchange = map[string]float64{}
		priceByExchange = append(priceByExchange, pricefromCurrentExchange)
		go p.getPriceFromExchange(pricefromCurrentExchange, exchangeSDK)
	}
	wg.Wait()

	// merge prices from all exchange sources
	for _, prices := range priceByExchange {
		for _, currency := range currencies {
			price, found := prices[currency]
			if found {
				priceListByCurrency[currency] = append(priceListByCurrency[currency], price)
			}
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

		// no data, skip
		if priceLen == 0 {
			continue
		}

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

	return medianPriceByCurrency
}

func (p *pricingImpl) getPriceFromExchange(pricefromCurrentExchange map[string]float64, exchangeSDK ExchangeSDK) {
	defer wg.Done()
	prices, err := exchangeSDK.GetPrices(currencies)
	if err != nil {
		// IMPROVEMENT: add logs and alarm
		return
	}
	for _, currency := range currencies {
		priceForCurrency, found := prices[currency]
		if !found {
			// IMPROVEMENT: logs error and alarm
			continue
		}
		pricefromCurrentExchange[currency] = priceForCurrency
	}
}
