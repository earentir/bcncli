package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

const apiURL = "https://bconomy.net/api/data"

// validateAPIKey exits if missing
func validateAPIKey() string {
	key := viper.GetString("apikey")
	if key == "" {
		fmt.Fprintln(os.Stderr, "Error: API key must be set via --apikey flag, config file, or env var BCONOMYAPI")
		os.Exit(1)
	}
	return key
}

// fetchData posts payload and returns bytes or error
func fetchData(payload map[string]interface{}, apiKey string) ([]byte, error) {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// FetchDataOrExit wraps fetchData
func FetchDataOrExit(payload map[string]interface{}) []byte {
	apiKey := validateAPIKey()
	data, err := fetchData(payload, apiKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching data: %v\n", err)
		os.Exit(1)
	}
	return data
}
