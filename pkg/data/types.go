package data_types

type MarketData struct {
	Ticker    string
	Timestamp int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int64
}