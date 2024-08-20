package clients

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	data_types "goquant/pkg/data"
)

type AlphaVantageClient struct {
	APIKey  string
	BaseURL string
}

// NewAlphaVantageClient creates a new AlphaVantageClient instance.
//
// Parameter apiKey is the API key used for authentication with Alpha Vantage.
// Returns a pointer to an AlphaVantageClient object.
func NewAlphaVantageClient(apiKey string) *AlphaVantageClient {

	return &AlphaVantageClient{
		APIKey:  apiKey,
		BaseURL: "https://www.alphavantage.co/query?",
	}
}

// FetchMinuteData fetches intraday minute data for a given symbol.
//
// Parameter symbol is the stock symbol to fetch data for.
// Parameter start and end are the start and end timestamps to filter data by.
// Parameter interval is the interval of the data (e.g. 1min, 5min, etc.).
// Returns a slice of MarketData and an error.
func (c *AlphaVantageClient) FetchMinuteData(symbol string, start, end int64, interval string) ([]data_types.MarketData, error) {
	// Prepare the URL
	url := fmt.Sprintf("%sfunction=TIME_SERIES_INTRADAY&symbol=%s&interval=%s&apikey=%s&datatype=csv&outputsize=full", c.BaseURL, symbol, interval, c.APIKey)
	fmt.Println(url)
	// Perform the request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	defer resp.Body.Close()

	// Parse the CSV response
	r := csv.NewReader(resp.Body)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	// Process the data and filter by the start and end time
	var data []data_types.MarketData
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		timestamp, _ := time.Parse("2006-01-02 15:04:05", record[0])
		if timestamp.Unix() < start || timestamp.Unix() > end {
			continue
		}
		open, _ := strconv.ParseFloat(record[1], 64)
		high, _ := strconv.ParseFloat(record[2], 64)
		low, _ := strconv.ParseFloat(record[3], 64)
		closePrice, _ := strconv.ParseFloat(record[4], 64)
		volume, _ := strconv.ParseInt(record[5], 10, 64)

		data = append(data, data_types.MarketData{
			Timestamp: timestamp.Unix(),
			Ticker:    symbol,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     closePrice,
			Volume:    volume,
		})
	}

	return data, nil
}
