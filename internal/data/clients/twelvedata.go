package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	data_types "goquant/pkg/data"
)

type TwelveDataClient struct {
	APIKey  string
	BaseURL string
}

// NewTwelveDataClient creates a new TwelveDataClient instance.
//
// Parameter apiKey is the API key used for authentication with Twelve Data.
// Returns a pointer to a TwelveDataClient object.
func NewTwelveDataClient(apiKey string) *TwelveDataClient {
	return &TwelveDataClient{
		APIKey:  apiKey,
		BaseURL: "https://api.twelvedata.com/time_series?",
	}
}

// FetchMinuteData fetches intraday minute data for a given symbol.
//
// Parameter symbol is the stock symbol to fetch data for.
// Parameter start and end are the start and end timestamps to filter data by.
// Parameter interval is the interval of the data (e.g., 1min, 5min, etc.).
// Returns a slice of MarketData and an error.
func (c *TwelveDataClient) FetchMinuteData(symbol string, start, end int64, interval string) ([]data_types.MarketData, error) {
	// Prepare the URL
	url := fmt.Sprintf("%ssymbol=%s&interval=%s&apikey=%s&start_date=%s&end_date=%s&format=JSON",
		c.BaseURL, symbol, interval, c.APIKey,
		time.Unix(start, 0).Format("2006-01-02 15:04:05"),
		time.Unix(end, 0).Format("2006-01-02 15:04:05"))

	// Perform the request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	body, error := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if error != nil {
		return nil, fmt.Errorf("error reading body: %v", error)
	}
	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("error reading body: %v", err)
	}

	// Parse the JSON response
	var result struct {
		Values []struct {
			Datetime string  `json:"datetime"`
			Open     float64 `json:"open,string"`
			High     float64 `json:"high,string"`
			Low      float64 `json:"low,string"`
			Close    float64 `json:"close,string"`
			Volume   int64   `json:"volume,string"`
		} `json:"values"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Process the data and filter by the start and end time
	var data []data_types.MarketData
	for _, record := range result.Values {
		timestamp, _ := time.Parse("2006-01-02 15:04:05", record.Datetime)
		if timestamp.Unix() < start || timestamp.Unix() > end {
			continue
		}
		data = append(data, data_types.MarketData{
			Timestamp: timestamp.Unix(),
			Ticker:    symbol,
			Open:      record.Open,
			High:      record.High,
			Low:       record.Low,
			Close:     record.Close,
			Volume:    record.Volume,
		})
	}

	return data, nil
}
