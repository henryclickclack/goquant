package strategies

import (
	backtest_types "goquant/pkg/backtest"
	"math"

	"github.com/go-gota/gota/dataframe"
)

// BollingerBandsReversionStrategy implements the Bollinger Bands reversion strategy
func BollingerBandsReversionStrategy(df dataframe.DataFrame) backtest_types.StrategyAction {
	period := 20            // Moving average period
	stdDevMultiplier := 2.0 // Standard deviation multiplier

	// Ensure we have enough data to calculate Bollinger Bands
	if df.Nrow() < period {
		return "Hold"
	}

	// Calculate the moving average and standard deviation for the closing prices
	prices := df.Col("Close").Float()
	movingAvg := movingAverage(prices, period)
	stdDev := standardDeviation(prices, movingAvg, period)

	// Calculate the upper and lower Bollinger Bands
	upperBand := make([]float64, len(movingAvg))
	lowerBand := make([]float64, len(movingAvg))
	for i := 0; i < len(movingAvg); i++ {
		upperBand[i] = movingAvg[i] + stdDevMultiplier*stdDev[i]
		lowerBand[i] = movingAvg[i] - stdDevMultiplier*stdDev[i]
	}

	// Get the current price and Bollinger Bands values
	currentPrice := prices[len(prices)-1]
	currentUpperBand := upperBand[len(upperBand)-1]
	currentLowerBand := lowerBand[len(lowerBand)-1]

	// Generate signals based on the price's position relative to the Bollinger Bands
	if currentPrice <= currentLowerBand {
		return "Buy"
	} else if currentPrice >= currentUpperBand {
		return "Sell"
	}

	return "Hold"
}

// standardDeviation calculates the standard deviation for a given period
func standardDeviation(data, movingAvg []float64, period int) []float64 {
	stdDev := make([]float64, len(data))

	for i := 0; i < len(data); i++ {
		if i+1 < period {
			stdDev[i] = 0 // Not enough data points yet
		} else {
			sum := 0.0
			for j := 0; j < period; j++ {
				sum += math.Pow(data[i-j]-movingAvg[i], 2)
			}
			stdDev[i] = math.Sqrt(sum / float64(period))
		}
	}

	return stdDev
}
