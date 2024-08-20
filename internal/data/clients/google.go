package clients

import (
	"encoding/csv"
	"fmt"
	data_types "goquant/pkg/data"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const googleFinanceURL = "https://finance.google.com/finance/historical?q="

type GoogleFinanceDataSource struct {
	Client *http.Client
}

// NewGoogleFinanceDataSource creates a new GoogleFinanceDataSource with a default HTTP client.
//
// No parameters.
// Returns a pointer to a GoogleFinanceDataSource object.
func NewGoogleFinanceDataSource() *GoogleFinanceDataSource {
	return &GoogleFinanceDataSource{
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Fetch retrieves market data from Google Finance.
//
// Parameters:
//
//	symbol: the stock symbol of the market data to be retrieved.
//	start: the start date of the market data in Unix timestamp.
//	end: the end date of the market data in Unix timestamp.
//
// Returns:
//
//	A slice of data_types.MarketData containing the retrieved market data.
//	An error if the data retrieval or parsing fails.
func (g *GoogleFinanceDataSource) Fetch(symbol string, start, end int64) ([]data_types.MarketData, error) {
	// Convert Unix timestamp to date format required by Google Finance (yyyy-MM-dd)
	startDate := time.Unix(start, 0).Format("2006-01-02")
	endDate := time.Unix(end, 0).Format("2006-01-02")

	url := fmt.Sprintf("%s%s&startdate=%s&enddate=%s&output=csv", googleFinanceURL, symbol, startDate, endDate)

	resp, err := g.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from Google Finance: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV data: %v", err)
	}

	marketData := convertGoogleFinanceData(symbol, records)
	return marketData, nil
}

// convertGoogleFinanceData converts Google Finance CSV response to MarketData slice.
//
// Parameters:
//
//	symbol: the stock symbol of the market data.
//	records: a 2D slice of strings representing the CSV data.
//
// Returns:
//
//	A slice of data_types.MarketData containing the converted market data.
func convertGoogleFinanceData(symbol string, records [][]string) []data_types.MarketData {
	var data []data_types.MarketData

	// Skip the first line (header)
	for _, record := range records[1:] {
		timestamp, _ := time.Parse("2-Jan-2006", record[0])
		open, _ := strconv.ParseFloat(record[1], 64)
		high, _ := strconv.ParseFloat(record[2], 64)
		low, _ := strconv.ParseFloat(record[3], 64)
		closePrice, _ := strconv.ParseFloat(record[4], 64)
		volume, _ := strconv.ParseInt(strings.ReplaceAll(record[5], ",", ""), 10, 64)

		data = append(data, data_types.MarketData{
			Ticker:    symbol,
			Timestamp: timestamp.Unix(),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     closePrice,
			Volume:    volume,
		})
	}
	return data
}
