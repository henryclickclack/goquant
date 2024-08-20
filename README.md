# goquant

Quant Library for Go

# GoQuant: A Quantitative Finance Library for Go

==============================================

GoQuant is a Go module designed to provide a comprehensive set of tools and libraries for quantitative finance. It aims to make it easy to work with financial data, perform calculations, and build trading strategies in the Go programming language.

## Features

-   **Data Retrieval**: GoQuant provides an interface to fetch historical market data from various sources, including Yahoo Finance, Google Finance, and IEX Cloud.
-   **Data Transformation**: Easily manipulate and transform financial data using GoQuant's built-in functions for data cleaning, filtering, and aggregation.
-   **Technical Indicators**: Calculate popular technical indicators such as Moving Averages, RSI, Bollinger Bands, and more.
-   **Strategy Backtesting**: Backtest trading strategies using GoQuant's built-in backtesting framework.

## Installation

To install GoQuant, run the following command:

```bash
go get github.com/goquant/goquant
```

## Examples

### Retrieving Historical Market Data

```go
package main

import (
	"fmt"
	"github.com/goquant/goquant/data"
)

func main() {
	// Create a new data client
	client := data.NewYahooFinanceDataSource()

	// Fetch historical data for Apple stock
	data, err := client.Fetch("AAPL", 0, 1643723900)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the data
	for _, row := range data {
		fmt.Println(row.Timestamp, row.Open, row.High, row.Low, row.Close)
	}
}
```

### Calculating Technical Indicators

```go
package main

import (
	"fmt"
	"github.com/goquant/goquant/indicators"
)

func main() {
	// Create a new data client
	client := data.NewYahooFinanceDataSource()

	// Fetch historical data for Apple stock
	data, err := client.Fetch("AAPL", 0, 1643723900)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Calculate the 50-day moving average
	ma := indicators.NewMovingAverage(50)
	for _, row := range data {
		ma.Update(row.Close)
		fmt.Println(row.Timestamp, ma.Value())
	}
}
```

### Backtesting a Trading Strategy

```go
package main

import (
	"fmt"
	"github.com/goquant/goquant/backtest"
)

func main() {
	// Create a new backtesting engine
	engine := backtest.NewEngine()

	// Define a simple trading strategy
	strategy := func(data []data.MarketData) {
		// Buy when the price is above the 50-day moving average
		if data[len(data)-1].Close > indicators.NewMovingAverage(50).Value() {
			engine.Buy(data[len(data)-1].Timestamp)
		}
	}

	// Backtest the strategy
	engine.Backtest(strategy, "AAPL", 0, 1643723900)

	// Print the results
	fmt.Println(engine.Results())
}
```

## Contributing

GoQuant is an open-source project and welcomes contributions from the community. If you're interested in contributing, please fork the repository and submit a pull request.

## License

GoQuant is licensed under the MIT License. See the LICENSE file for details.
