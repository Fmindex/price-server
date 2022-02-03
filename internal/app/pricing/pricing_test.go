package pricing

import (
	"errors"
	"testing"

	"github.com/Fmindex/price-server/internal/app/pricing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	nilError error
)

// IMPROVEMENT: add test cases
func TestPricing(t *testing.T) {
	var testCases = []struct {
		scenario                string
		exchange1ReturnTestcase map[string]float64
		exchange1ErrorTestcase  error
		exchange2ReturnTestcase map[string]float64
		exchange2ErrorTestcase  error
		output                  map[string]string
	}{
		{
			scenario:                "normal request",
			exchange1ReturnTestcase: map[string]float64{"LUNA": 2.0, "BTC": 1.0},
			exchange1ErrorTestcase:  nilError,
			exchange2ReturnTestcase: map[string]float64{"LUNA": 8.0},
			exchange2ErrorTestcase:  nilError,
			output:                  map[string]string{"LUNA": "5.000000", "BTC": "1.000000"},
		},
		{
			scenario:                "exchange sdk 2 failed",
			exchange1ReturnTestcase: map[string]float64{"LUNA": 2.0},
			exchange1ErrorTestcase:  nilError,
			exchange2ReturnTestcase: map[string]float64{},
			exchange2ErrorTestcase:  errors.New("currency is not supported"),
			output:                  map[string]string{"LUNA": "2.000000"},
		},
	}

	assert := assert.New(t)
	mockedExchangeSDK1 := mocks.ExchangeSDK{}
	mockedExchangeSDK2 := mocks.ExchangeSDK{}
	mockedExchangeSDKs := []ExchangeSDK{&mockedExchangeSDK1, &mockedExchangeSDK2}
	app := NewPricing(mockedExchangeSDKs)

	for _, c := range testCases {

		// set SDK1 returns
		mockedExchangeSDK1.On("GetPrices", mock.Anything).Return(
			func(currencies []string) map[string]float64 {
				return c.exchange1ReturnTestcase
			},
			func(currencies []string) error {
				return c.exchange1ErrorTestcase
			})

		// set SDK2 returns
		mockedExchangeSDK2.On("GetPrices", mock.Anything).Return(
			func(currencies []string) map[string]float64 {
				return c.exchange2ReturnTestcase
			},
			func(currencies []string) error {
				return c.exchange2ErrorTestcase
			})

		res := app.getLatestPriceImpl()
		assert.Equal(c.output, res, c.scenario)
	}
}
