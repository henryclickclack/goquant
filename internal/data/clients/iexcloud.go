package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	data_types "goquant/pkg/data"
)

type IEXCloudClient struct {
	APIKey  string
	BaseURL string
}

// NewIEXCloudClient creates a new IEXCloudClient instance.
//
// Parameter apiKey is the API key used for authentication with IEX Cloud.
// Returns a pointer to an IEXCloudClient object.
func NewIEXCloudClient(apiKey string) *IEXCloudClient {
	return &IEXCloudClient{
		APIKey:  apiKey,
		BaseURL: "https://cloud.iexapis.com/stable/stock/",
	}
}

// FetchMinuteData fetches intraday minute data for a given symbol.
//
// Parameter symbol is the stock symbol to fetch data for.
// Parameter start and end are the start and end timestamps to filter data by.
// Parameter interval is ignored since IEX Cloud provides data at 1-minute intervals.
// Returns a slice of MarketData and an error.
func (c *IEXCloudClient) FetchMinuteData(symbol string, start, end int64, interval string) ([]data_types.MarketData, error) {
	// Prepare the URL
	url := fmt.Sprintf("%s%s/intraday-prices?token=%s", c.BaseURL, symbol, c.APIKey)

	// Perform the request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %v", err)
	}

	// Parse the JSON response
	var result []struct {
		Minute string  `json:"minute"`
		Open   float64 `json:"open"`
		High   float64 `json:"high"`
		Low    float64 `json:"low"`
		Close  float64 `json:"close"`
		Volume int64   `json:"volume"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Process the data and filter by the start and end time
	var data []data_types.MarketData
	for _, record := range result {
		timestamp, _ := time.Parse("15:04", record.Minute)
		fullTimestamp := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), timestamp.Hour(), timestamp.Minute(), 0, 0, time.UTC)

		if fullTimestamp.Unix() < start || fullTimestamp.Unix() > end {
			continue
		}

		data = append(data, data_types.MarketData{
			Timestamp: fullTimestamp.Unix(),
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
