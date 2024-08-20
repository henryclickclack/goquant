package config

import (
	"os"
)

type Config struct {
    APIKey     string
    DataSource string
}

func LoadConfig() Config {
    return Config{
        APIKey:     os.Getenv("API_KEY"),
        DataSource: os.Getenv("DATA_SOURCE"),
    }
}