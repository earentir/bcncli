package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

const apiURL = "https://bconomy.net/api/data"

// validateAPIKey exits if missing
func ValidateAPIKey() string {
	key := viper.GetString("apikey")
	if key == "" {
		fmt.Fprintln(os.Stderr, "Error: API key must be set via --apikey flag, config file, or env var BCONOMYAPI")
		os.Exit(1)
	}
	return key
}

// parseID converts arg to int
func ParseID(arg string) int {
	id, err := strconv.Atoi(arg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid ID: %s\n", arg)
		os.Exit(1)
	}
	return id
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
	apiKey := ValidateAPIKey()
	data, err := fetchData(payload, apiKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching data: %v\n", err)
		os.Exit(1)
	}
	return data
}

// PrintJSON pretty-prints
func PrintJSON(data []byte) {
	var pretty bytes.Buffer
	json.Indent(&pretty, data, "", "  ")
	fmt.Println(pretty.String())
}

// EpochToISO8601 converts milliseconds to ISO 8601 format
func EpochToISO8601(ms int64) string {
	if ms <= 0 {
		return "-"
	}
	// convert ms to nanoseconds, cast to int64
	nanos := int64(time.Duration(ms) * time.Millisecond)
	t := time.Unix(0, nanos).UTC()
	return t.Format(time.RFC3339)
}
