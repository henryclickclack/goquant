package main

import (
	"fmt"
	backtest "goquant/internal/backtesting"
	"goquant/internal/data/clients"
	"goquant/internal/data/storage"
	"goquant/internal/strategies"

	"time"
)

func main() {
	// Initialize a DataSource (Yahoo Finance in this case)
	//dataSource := clients.NewAlphaVantageClient("9EQCLVN4UJLBPQU9")
	dataSource := clients.NewTwelveDataClient("228c3167a61f4916b63323025b5f2165")
	//dataSource := clients.NewYahooFinanceDataSource()
	storage := storage.NewInMemoryStorage()

	// Set the time range (example: last 30 days)
	end := time.Now()
	start := end.Add(-time.Hour * 24 * 30)

	// Fetch data for a symbol (e.g., "AAPL") using the DataSource interface
	marketData, err := dataSource.FetchMinuteData("AAPL", start.Unix(), end.Unix(), "15min")
	//marketData, err := dataSource.Fetch("TSLA", start, end)
	fmt.Println(marketData)
	//marketData, err := dataSource.Fetch("FIVE", start, end)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}

	// Save data to storage
	err = storage.Save(marketData)
	if err != nil {
		fmt.Printf("Error saving data: %v\n", err)
		return
	}

	// Convert saved data to DataFrame
	df := storage.ToDataFrame(marketData)

	fmt.Println(df)

	initialInvest := 10000.0
	// Define individual strategies
	markovStrategy := strategies.NewMarkovChainStrategy(2)
	markovStrategy.Build(df) // Build the Markov model

	movingAverageStrategy := strategies.MovingAverageCrossoverStrategy

	// Define an ensemble strategy using the individual strategies
	ensemble := strategies.NewEnsembleStrategy([]strategies.StrategyFunc{
		markovStrategy.Run,
		movingAverageStrategy,
	}, []float64{0.5, 0.5}) // Equal weights

	// Run the backtest using the ensemble strategy
	result, err := backtest.Backtest(df, ensemble.Run, time.Minute*15, initialInvest)
	if err != nil {
		fmt.Printf("Backtest error: %v\n", err)
		return
	}
	fmt.Println("Sell count: ", result.SellCount)
	fmt.Println("Buy count: ", result.BuyCount)
	fmt.Println("Hold count: ", result.HoldCount)

	fmt.Println("Total profit/loss: ", result.TotalProfitLoss)
	fmt.Println("Gain strategy: ", result.GainStrategy)

	fmt.Println("Max up: ", result.MaxUp)
	fmt.Println("Max down: ", result.MaxDown)
	fmt.Println("Gain market: ", result.GainMarket)

	fmt.Println("Gain vs. market: ", result.GainVsMarket)
}
