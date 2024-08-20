package storage

import (
	"fmt"
	data_types "goquant/pkg/data"

	"github.com/go-gota/gota/dataframe"
)

type InMemoryStorage struct {
    store map[string][]data_types.MarketData
}

// NewInMemoryStorage creates a new InMemoryStorage instance.
//
// No parameters.
// Returns a pointer to an InMemoryStorage object.
func NewInMemoryStorage() *InMemoryStorage {
    return &InMemoryStorage{
        store: make(map[string][]data_types.MarketData), 
    }
}


// Save stores the market data in memory.
//
// It takes a slice of data_types.MarketData as a parameter.
// It appends each element of the slice to the corresponding ticker in the store.
// It returns an error if any.
func (s *InMemoryStorage) Save(data []data_types.MarketData) error {
    for _, d := range data {
        s.store[d.Ticker] = append(s.store[d.Ticker], d)
    }
    return nil
}


// Load retrieves the market data from memory for a given symbol within a specified time range.
//
// Parameters:
//   symbol (string): The stock symbol to retrieve data for.
//   start (int64): The start date of the time range to retrieve data for.
//   end (int64): The end date of the time range to retrieve data for.
// Returns:
//   A slice of data_types.MarketData containing the retrieved market data.
//   An error if the data retrieval fails.
func (s *InMemoryStorage) Load(symbol string, start, end int64) ([]data_types.MarketData, error) {
    data, ok := s.store[symbol]
    if !ok {
        return nil, fmt.Errorf("no data found for symbol: %s", symbol)
    }

    var filteredData []data_types.MarketData
    for _, d := range data {
        if d.Timestamp >= start && d.Timestamp <= end {
            filteredData = append(filteredData, d)
        }
    }
    return filteredData, nil
}


// ToDataFrame converts a slice of MarketData to a DataFrame.
//
// Parameters:
//   data ([]data_types.MarketData): The market data to be converted.
// Returns:
//   dataframe.DataFrame: The converted DataFrame.
func (s *InMemoryStorage) ToDataFrame(data []data_types.MarketData) dataframe.DataFrame {
    var records []map[string]interface{}

    for _, d := range data {
        record := map[string]interface{}{
            "Ticker":    d.Ticker,
            "Timestamp": d.Timestamp,
            "Open":      d.Open,
            "High":      d.High,
            "Low":       d.Low,
            "Close":     d.Close,
            "Volume":    d.Volume,
        }
        records = append(records, record)
    }

    df := dataframe.LoadMaps(records)
    return df
}
