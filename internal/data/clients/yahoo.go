package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	data_types "goquant/pkg/data"
)

const yahooFinanceURL = "https://query1.finance.yahoo.com/v8/finance/chart/"


type YahooFinanceDataSource struct {
    Client *http.Client
}


// NewYahooFinanceDataSource creates a new YahooFinanceDataSource with a default HTTP client.
//
// No parameters.
// Returns a pointer to a YahooFinanceDataSource object.
func NewYahooFinanceDataSource() *YahooFinanceDataSource {
    return &YahooFinanceDataSource{
        Client: &http.Client{Timeout: 10 * time.Second},
    }
}


// Fetch retrieves market data_types from Yahoo Finance
//
// Parameters:
//   symbol (string): The stock symbol to fetch data for.
//   start (int64): The start date of the time range to fetch data for.
//   end (int64): The end date of the time range to fetch data for.
// Returns:
//   A slice of data_types.MarketData and an error.
func (y *YahooFinanceDataSource) Fetch(symbol string, start, end int64) ([]data_types.MarketData, error) {
    url := fmt.Sprintf("%s%s?period1=%d&period2=%d&interval=1d", yahooFinanceURL, symbol, start, end)

    resp, err := y.Client.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch data_types from Yahoo Finance: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
    }

    var result yahooFinanceResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %v", err)
    }

    marketData := convertToMarketData(symbol, result)
    return marketData, nil
}


// convertToMarketData Converts a Yahoo Finance response to a slice of MarketData.
//
// Parameters:
//   symbol (string): The stock symbol of the MarketData.
//   response (yahooFinanceResponse): The Yahoo Finance response to be converted.
// Returns:
//   A slice of data_types.MarketData.
func convertToMarketData(symbol string, response yahooFinanceResponse) []data_types.MarketData {
    result := response.Chart.Result[0]
    timestamps := result.Timestamp
    quotes := result.Indicators.Quote[0]

    var dt []data_types.MarketData
    for i, timestamp := range timestamps {
        dt = append(dt, data_types.MarketData{
            Ticker:    symbol,
            Timestamp: timestamp,
            Open:      quotes.Open[i],
            High:      quotes.High[i],
            Low:       quotes.Low[i],
            Close:     quotes.Close[i],
            Volume:    quotes.Volume[i],
        })
    }
    return dt
}


type yahooFinanceResponse struct {
    Chart struct {
        Result []struct {
            Timestamp  []int64 `json:"timestamp"`
            Indicators struct {
                Quote []struct {
                    Open   []float64 `json:"open"`
                    High   []float64 `json:"high"`
                    Low    []float64 `json:"low"`
                    Close  []float64 `json:"close"`
                    Volume []int64   `json:"volume"`
                } `json:"quote"`
            } `json:"indicators"`
        } `json:"result"`
    } `json:"chart"`
}
