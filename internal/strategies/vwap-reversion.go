package strategies

import (
	backtest_types "goquant/pkg/backtest"

	"github.com/go-gota/gota/dataframe"
)

// VWAPReversionStrategy implements the VWAP reversion strategy.
//
// It takes a dataframe as input and returns a backtest_types.StrategyAction.
// The strategy calculates the VWAP (Volume Weighted Average Price) for the given data and generates signals based on the price deviation from the VWAP.
// If the price deviation is less than -0.01, it returns "Buy". If the price deviation is greater than 0.01, it returns "Sell". Otherwise, it returns "Hold".
func VWAPReversionStrategy(df dataframe.DataFrame) backtest_types.StrategyAction {
	// Ensure we have enough data to calculate VWAP
	if df.Nrow() < 1 {
		return "Hold"
	}

	// Calculate VWAP
	vwap := calculateVWAP(df.Col("Close").Float(), df.Col("Volume").Float())

	// Get the current price and VWAP
	currentPrice := df.Col("Close").Float()[df.Nrow()-1]
	currentVWAP := vwap[len(vwap)-1]

	// Calculate the price deviation from VWAP
	deviation := (currentPrice - currentVWAP) / currentVWAP

	// Generate signals based on price deviation from VWAP
	if deviation < -0.01 { // Buy if price is 1% below VWAP
		return "Buy"
	} else if deviation > 0.01 { // Sell if price is 1% above VWAP
		return "Sell"
	}

	return "Hold"
}

// calculateVWAP calculates the Volume Weighted Average Price (VWAP) for the given price and volume data.
//
// It takes two parameters: prices and volumes, both of which are slices of float64 representing the price and volume data respectively.
// It returns a slice of float64 representing the calculated VWAP for each data point.
func calculateVWAP(prices, volumes []float64) []float64 {
	vwap := make([]float64, len(prices))
	cumulativePriceVolume := 0.0
	cumulativeVolume := 0.0

	for i := 0; i < len(prices); i++ {
		cumulativePriceVolume += prices[i] * volumes[i]
		cumulativeVolume += volumes[i]
		vwap[i] = cumulativePriceVolume / cumulativeVolume
	}

	return vwap
}
