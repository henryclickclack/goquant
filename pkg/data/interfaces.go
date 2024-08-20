package data_types

import "github.com/go-gota/gota/dataframe"

type DataSource interface {
	Fetch(symbol string, start, end int64) ([]MarketData, error)
}

type DataStorage interface {
    Save(data []MarketData) error
    Load(symbol string, start, end int64) ([]MarketData, error)
    ToDataFrame(data []MarketData) dataframe.DataFrame
}