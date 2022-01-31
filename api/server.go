package server

import (
	"net/http"

	"github.com/Fmindex/price-server/internal/app/pricing"
	"github.com/Fmindex/price-server/internal/pkg/binance"
)

func Run() {

	// Init
	// SDKs
	binanceSDK := binance.NewBinanceSDK()
	exchangeSDK := []pricing.ExchangeSDK{binanceSDK}
	// handler
	pricingHandler := pricing.NewPricing(exchangeSDK)

	http.HandleFunc("/latest", pricingHandler.GetLatestPrice)

	http.ListenAndServe(":8888", nil)
}
