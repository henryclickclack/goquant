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
	//dataSource := clients.NewTwelveDataClient("228c3167a61f4916b63323025b5f2165")
	dataSource := clients.NewYahooFinanceDataSource()
	storage := storage.NewInMemoryStorage()

	// Set the time range (example: last 30 days)
	end := time.Now().Unix()
	start := end - 365*24*3600

	// Fetch data for a symbol (e.g., "AAPL") using the DataSource interface
	//marketData, err := dataSource.FetchMinuteData("AAPL", start, end, "1min")
	marketData, err := dataSource.Fetch("TSLA", start, end)
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
	markovStrategy := strategies.NewMarkovChainStrategy()
	markovStrategy.Build(df)
	// Run backtest
	result, err := backtest.Backtest(df, markovStrategy.Run, time.Hour*24, initialInvest)
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