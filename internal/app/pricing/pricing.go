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
	GetPrices(symbols []string) (map[string]float64, error)
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
	symbols = []string{"BTC", "ETH", "LUNA"}
	wg      sync.WaitGroup
)

type GetLatestPriceResponse struct {
	Prices map[string]string `json:"prices"`
}

// Business Logic
// GetLatestPrice is to get the latest price of symbols from each SDKs
func (p *pricingImpl) GetLatestPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get price
	medianPriceBySymbol := p.getLatestPriceImpl()

	// api return
	fmt.Printf("/latest: return %v\n", medianPriceBySymbol)
	json.NewEncoder(w).Encode(GetLatestPriceResponse{
		Prices: medianPriceBySymbol,
	})
}

func (p *pricingImpl) getLatestPriceImpl() map[string]string {

	var priceListBySymbol = map[string][]float64{}
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
		for _, symbol := range symbols {
			price, found := prices[symbol]
			if found {
				priceListBySymbol[symbol] = append(priceListBySymbol[symbol], price)
			}
		}
	}

	// find median from each symbols
	// sort array
	for _, symbol := range symbols {
		sort.Float64s(priceListBySymbol[symbol])
	}
	// find median
	var medianPriceBySymbol = map[string]string{}
	for _, symbol := range symbols {
		priceLen := len(priceListBySymbol[symbol])

		// no data, skip
		if priceLen == 0 {
			continue
		}

		var price float64
		if priceLen%2 == 0 {
			// even, average of the middles
			price = (priceListBySymbol[symbol][priceLen/2-1] + priceListBySymbol[symbol][priceLen/2]) / 2.0
		} else {
			// odd, use the middle one
			price = priceListBySymbol[symbol][priceLen/2]
		}
		medianPriceBySymbol[symbol] = fmt.Sprintf("%f", price)
	}

	return medianPriceBySymbol
}

func (p *pricingImpl) getPriceFromExchange(pricefromCurrentExchange map[string]float64, exchangeSDK ExchangeSDK) {
	defer wg.Done()
	prices, err := exchangeSDK.GetPrices(symbols)
	if err != nil {
		// IMPROVEMENT: add logs and alarm
		return
	}
	for _, symbol := range symbols {
		priceForSymbol, found := prices[symbol]
		if !found {
			// IMPROVEMENT: logs error and alarm
			continue
		}
		pricefromCurrentExchange[symbol] = priceForSymbol
	}
}
