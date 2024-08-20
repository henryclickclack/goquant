package backtest_types

import "github.com/go-gota/gota/dataframe"

//TODO"Buy" | "Sell" | "Hold"
type StrategyAction string

func (s StrategyAction) Valid() bool {
	switch s {
	case "Buy", "Sell", "Hold":
		return true
	default:
		return false
	}
}

type StrategyFunction func(df dataframe.DataFrame) StrategyAction

type BacktestResult struct {
	TotalProfitLoss float64
	MaxUp           float64
	MaxDown         float64
	TradeLog        dataframe.DataFrame
	BuyCount        int
	SellCount       int
	HoldCount       int
	TotalCount      int
	GainMarket      float64
	GainStrategy    float64
	GainVsMarket    float64
}
