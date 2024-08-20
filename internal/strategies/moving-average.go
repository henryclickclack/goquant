package strategies

import (
	"fmt"
	backtest_types "goquant/pkg/backtest"

	"github.com/go-gota/gota/dataframe"
)

// MovingAverageCrossoverStrategy implements a moving average crossover strategy.
//
// Parameters:
// - df (dataframe.DataFrame): the input dataframe containing the financial data.
// Returns:
// - backtest_types.StrategyAction: the trading signal ("Buy", "Sell", or "Hold") based on the moving average crossover.
func MovingAverageCrossoverStrategy(df dataframe.DataFrame) backtest_types.StrategyAction {
	shortWindow := 5 // Short-term moving average window
	longWindow := 20 // Long-term moving average window

	// Ensure we have enough data to calculate both moving averages
	if df.Nrow() < longWindow {
		fmt.Printf("not enough data to calculate both moving averages\n")
		return "Hold"
	}

	// Calculate the short-term moving average
	shortMA := movingAverage(df.Col("Close").Float(), shortWindow)

	// Calculate the long-term moving average
	longMA := movingAverage(df.Col("Close").Float(), longWindow)

	// Get the last short and long moving average values
	currentShortMA := shortMA[len(shortMA)-1]
	currentLongMA := longMA[len(longMA)-1]

	// Get the previous short and long moving average values (to detect crossover)
	prevShortMA := shortMA[len(shortMA)-2]
	prevLongMA := longMA[len(longMA)-2]

	// Generate trading signals based on the crossover
	if prevShortMA <= prevLongMA && currentShortMA > currentLongMA {
		return "Buy"
	} else if prevShortMA >= prevLongMA && currentShortMA < currentLongMA {
		return "Sell"
	}

	return "Hold"
}

// movingAverage calculates the moving average of a given dataset.
//
// Parameters:
// - data (slice of float64): the input dataset.
// - window (int): the size of the moving average window.
// Returns:
// - slice of float64: the moving average values.
func movingAverage(data []float64, window int) []float64 {
	movingAvg := make([]float64, len(data))

	for i := 0; i < len(data); i++ {
		if i+1 < window {
			movingAvg[i] = 0 // Not enough data points yet, fill with 0 or NaN
		} else {
			sum := 0.0
			for j := 0; j < window; j++ {
				sum += data[i-j]
			}
			movingAvg[i] = sum / float64(window)
		}
	}

	return movingAvg
}
