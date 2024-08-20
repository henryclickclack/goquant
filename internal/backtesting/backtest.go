package backtest

import (
	"errors"
	"fmt"
	backtest_types "goquant/pkg/backtest"
	"slices"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

// StrategyFunction is a type that represents a strategy function

// Backtest runs a backtesting simulation on a given dataframe using a specified strategy function at a specified interval.
//
// Parameters:
//
//	df (dataframe.DataFrame): The input dataframe containing the financial data.
//	strategy (backtest_types.StrategyFunction): The strategy function to be applied to the dataframe.
//	interval (time.Duration): The interval at which the strategy function is applied.
//	initialInvest (float64): The initial investment amount.
//
// Returns:
//
//	backtest_types.BacktestResult: The result of the backtesting simulation.
//	error: Any error that occurred during the simulation.
func Backtest(df dataframe.DataFrame, strategy backtest_types.StrategyFunction, interval time.Duration, initialInvest float64) (backtest_types.BacktestResult, error) {
	// Ensure the dataframe has a "Timestamp" column
	if !slices.Contains(df.Names(), "Timestamp") {
		return backtest_types.BacktestResult{}, errors.New("dataframe must have a 'Timestamp' column")
	}

	// Check if the dataframe's interval is fine-grained enough
	if !isIntervalFineEnough(df, interval) {
		return backtest_types.BacktestResult{}, fmt.Errorf("data interval is not fine-grained enough for the specified interval: %v", interval)
	}

	// Initialize variables for tracking performance
	totalProfitLoss := 0.0
	maxUp := 0.0
	maxDown := 0.0
	currentInvest := initialInvest
	tradeResults := []map[string]interface{}{}

	timestamps := df.Col("Timestamp").Float()
	for i := 0; i < len(timestamps); i++ {
		if i > 0 && timeBetween(timestamps[i-1], timestamps[i]) < interval {
			fmt.Println(timeBetween(timestamps[i-1], timestamps[i]), "    <   ", interval)
			fmt.Println("Skipping interval")
			continue
		}
		// print the i
		fmt.Printf("i: %d\n", i)
		// Apply the strategy to the current row
		//set subset to 0:i
		subset := make([]int, i)
		for c := 0; c < i; c++ {
			subset[c] = c
		}
		action := strategy(df.Subset(subset))
		openPrice := df.Col("Open").Float()[i]
		closePrice := df.Col("Close").Float()[i]

		profitLoss := 0.0

		// Calculate profit or loss based on the action
		if action == "Buy" {
			profitLoss = (closePrice - openPrice) / openPrice * currentInvest
		} else if action == "Sell" {
			profitLoss = (openPrice - closePrice) / openPrice * currentInvest
		}

		// Update total profit/loss and current investment
		totalProfitLoss += profitLoss
		currentInvest += profitLoss

		// Ensure investment doesn't drop below zero
		if currentInvest < 0 {
			currentInvest = 0
		}

		// Track max up and max down
		if totalProfitLoss > maxUp {
			maxUp = totalProfitLoss
		}
		if totalProfitLoss < maxDown {
			maxDown = totalProfitLoss
		}

		// Record the trade result
		tradeResults = append(tradeResults, map[string]interface{}{
			"Timestamp":       time.Unix(int64(timestamps[i]), 0).Format(time.RFC3339),
			"Action":          action,
			"OpenPrice":       openPrice,
			"ClosePrice":      closePrice,
			"ProfitLoss":      profitLoss,
			"TotalProfitLoss": totalProfitLoss,
			"CurrentInvest":   currentInvest,
		})

		// Print the action (for debugging or logging purposes)
		fmt.Printf("Time: %v, Action: %s, Profit/Loss: %.2f, Total P/L: %.2f, Current Investment: %.2f\n",
			time.Unix(int64(timestamps[i]), 0), action, profitLoss, totalProfitLoss, currentInvest)
	}

	// Convert tradeResults to a dataframe
	tradeLogDF := dataframe.LoadMaps(tradeResults)

	// Calculate action counts
	buyCount := tradeLogDF.Filter(dataframe.F{Colname: "Action", Comparator: series.Eq, Comparando: "Buy"}).Nrow()
	sellCount := tradeLogDF.Filter(dataframe.F{Colname: "Action", Comparator: series.Eq, Comparando: "Sell"}).Nrow()
	holdCount := tradeLogDF.Filter(dataframe.F{Colname: "Action", Comparator: series.Eq, Comparando: "Hold"}).Nrow()
	//TotalPNLMarket := marketData[len(marketData)-1].Close - marketData[0].Open
	//gainMarket := marketData[len(marketData)-1].Close/marketData[0].Open - 1
	gainMarket := (tradeLogDF.Col("ClosePrice").Float()[len(tradeLogDF.Col("ClosePrice").Float())-1] - tradeLogDF.Col("OpenPrice").Float()[0]) / tradeLogDF.Col("OpenPrice").Float()[0]
	gainStrategy := totalProfitLoss / initialInvest
	gainVsMarket := gainStrategy - gainMarket

	// Return the backtest results
	return backtest_types.BacktestResult{
		TotalProfitLoss: totalProfitLoss,
		MaxUp:           maxUp,
		MaxDown:         maxDown,
		TradeLog:        tradeLogDF,
		BuyCount:        buyCount,
		SellCount:       sellCount,
		HoldCount:       holdCount,
		GainMarket:      gainMarket,
		GainStrategy:    gainStrategy,
		GainVsMarket:    gainVsMarket,
	}, nil
}

// isIntervalFineEnough checks if the dataframe's interval is fine-grained enough for the desired backtest interval.
//
// Parameters:
// - df: the dataframe to check.
// - requiredInterval: the desired interval.
//
// Returns:
// - true if the interval is fine-grained enough, false otherwise.
func isIntervalFineEnough(df dataframe.DataFrame, requiredInterval time.Duration) bool {
	timestamps := df.Col("Timestamp").Float()
	for i := 1; i < len(timestamps); i++ {
		if timeBetween(timestamps[i-1], timestamps[i]) <= requiredInterval {
			return true
		}
	}
	return false
}

// timeBetween calculates the time difference between two timestamps.
//
// Parameters:
// - timestamp1: the first timestamp in float64.
// - timestamp2: the second timestamp in float64.
//
// Returns:
// - the time.Duration representing the time difference between the two timestamps.
func timeBetween(timestamp1, timestamp2 float64) time.Duration {
	if timestamp1 > timestamp2 {
		return time.Duration(timestamp1-timestamp2) * time.Second
	}
	return time.Duration(timestamp2-timestamp1) * time.Second
}
