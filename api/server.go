package api

import (
	"fmt"
	"net/http"

	"github.com/Fmindex/price-server/internal/app/pricing"
	"github.com/Fmindex/price-server/internal/pkg/binance"
	"github.com/Fmindex/price-server/internal/pkg/coingecko"
)

func Run() {

	// Init
	// SDKs
	binanceSDK := binance.NewBinanceSDK()
	coingeckoSDK := coingecko.NewCoingeckoSDK()
	exchangeSDK := []pricing.ExchangeSDK{binanceSDK, coingeckoSDK}
	// handler
	pricingHandler := pricing.NewPricing(exchangeSDK)

	// register handler
	http.HandleFunc("/latest", pricingHandler.GetLatestPrice)

	// run http server
	fmt.Println("Price server is running")
	http.ListenAndServe(":8888", nil)
}
