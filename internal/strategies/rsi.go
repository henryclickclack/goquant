package strategies

import (
	"fmt"
	backtest_types "goquant/pkg/backtest"

	"github.com/go-gota/gota/dataframe"
)

// RSIStrategy generates a trading signal based on the Relative Strength Index (RSI) of a given DataFrame.
//
// Parameters:
// - df (dataframe.DataFrame): A DataFrame containing the financial data to be analyzed.
// Returns:
// - backtest_types.StrategyAction: A trading signal indicating whether to "Buy", "Sell", or "Hold".
func RSIStrategy(df dataframe.DataFrame) backtest_types.StrategyAction {
	period := 2 // Standard RSI period

	// Ensure we have enough data to calculate RSI
	fmt.Printf("Num of rows in dataframe: %d\n", df.Nrow())
	if df.Nrow()-1 < period {
		fmt.Printf("not enough data to calculate RSI\n")
		return "Hold"
	}

	// Calculate the RSI
	rsi := calculateRSI(df.Col("Close").Float(), period)
	// Get the current RSI value
	currentRSI := rsi[len(rsi)-1]
	//
	// Generate signals based on RSI levels
	if currentRSI < 30 {
		return "Buy"
	} else if currentRSI > 70 {
		return "Sell"
	}

	return "Hold"
}

// calculateRSI calculates the Relative Strength Index (RSI) for a given period
func calculateRSI(prices []float64, period int) []float64 {
	delta := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		delta[i-1] = prices[i] - prices[i-1]
	}

	gain := make([]float64, len(delta))
	loss := make([]float64, len(delta))
	for i, d := range delta {
		if d > 0 {
			gain[i] = d
		} else {
			loss[i] = -d
		}
	}

	avgGain := make([]float64, len(gain))
	avgLoss := make([]float64, len(loss))

	// Initial average gain and loss
	avgGain[period-1] = sum(gain[:period]) / float64(period)
	avgLoss[period-1] = sum(loss[:period]) / float64(period)

	// Calculate the rest of the average gains and losses
	for i := period; i < len(gain); i++ {
		avgGain[i] = (avgGain[i-1]*(float64(period)-1) + gain[i]) / float64(period)
		avgLoss[i] = (avgLoss[i-1]*(float64(period)-1) + loss[i]) / float64(period)
	}

	rsi := make([]float64, len(prices))
	for i := period - 1; i < len(avgGain); i++ {
		if avgLoss[i] == 0 {
			rsi[i+1] = 100
		} else {
			rs := avgGain[i] / avgLoss[i]
			rsi[i+1] = 100 - (100 / (1 + rs))
		}
	}

	return rsi
}

// sum calculates the sum of a slice of float64 values
func sum(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}
