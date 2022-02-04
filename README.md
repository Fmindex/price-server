# price-server

### Overview

An HTTP API sevice exposing API to get a price of configured symbols from each sources as configured.

#### Key concepts

- The part of code which connects to Exchange (e.g. Coingecko, Binance) are separated to be SDK. So, we can implement the business/parsing logic in each SDK independently. Only need to keep the interface the same as in the standard (has GetPrices method). The implemented SDK will be connected to main business logic using dependencies injection so that we can inject mocked SDK in order to test the main business logic.

- The symbols are also configurable. no need to update any code logic if only need new symbols.

- The API calls to each SDK are called concurrently making logic faster. All the concurrent calls will be wait by WaitGroups until all the calls are finished. After that, we can process to the next logic so we are sure all the wanted prices are there.

### Quick Start

1. Build: `make build`
2. Run: `make run`
3. (Alternative) Build and run: `make build_and_run`

### How to add new exchange price source

1. Add new package into the `internal/pkg/`
2. Implement logic to query API there expose the logic with method

   ```go
   GetPrices([]string) (map[string]float64, error)
   ```

3. In the `api/server.go`, add logic to init new SDK and inject to pricing service

### How to add new currency

1. In the `internal/app/pricing.go`, add the new symbols to `symbols`
2. In each SDK, add new symbols to `symbolMapper`

### Potential improvement plans

1. Add logging service to keep the error, warning and info logs.
2. If we have logs, we can have the alarm to alert us when the API call's error rate is too high. Or some logic errors are too high.
3. Can add more test cases and add code coverage metric to measure how coverage our test is.
4. The list of symbols can be maintain in some configuration outside of code logic avoiding code deployment if we only need to add new symbol.
5. the symbols can also become Enums helping us know if there is some type on symbol.
6. Can add some CI/CD automating our changes deployment.
7. Can dockerize it to make local testing portable for others.
8. If some source may charge us from too high API calls, we can implement API rate limit.
