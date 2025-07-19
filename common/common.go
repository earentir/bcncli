// Package common provides common utilities and types used across the application.
package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"bcncli/client"
)

// Item represents an entry from itemid.json, with every field included.
type Item struct {
	Name        string        `json:"name"`
	Emoji       string        `json:"emoji"`
	IDName      string        `json:"idName"`
	Uncraftable bool          `json:"uncraftable"`
	Attributes  []string      `json:"attributes"`
	LootSources []string      `json:"lootSources"`
	Recipe      []interface{} `json:"recipe"`
	ID          int           `json:"id"`
	FlatID      string        `json:"flatId"`
	Cost        int64         `json:"cost"`
	UsedToCraft []interface{} `json:"usedToCraft"`
	ImageURL    string        `json:"imageUrl"`
}

// FormatPrice formats n either with units (K, M, B, T) or, if you pass
// useNum=true, as a plain integer with space separators.
//
//	formatNumber(1230000)        == "1.23M"
//	formatNumber(1000000)        == "1M"
//	formatNumber(200000000)      == "200M"
//	formatNumber(123000000000)   == "123B"
//	formatNumber(123000, true)   == "123 000"
func FormatPrice(n int64, useNum ...bool) string {
	num := false
	if len(useNum) > 0 && useNum[0] {
		num = true
	}

	// Plain number with spaces
	if num {
		s := strconv.FormatInt(n, 10)
		neg := strings.HasPrefix(s, "-")
		if neg {
			s = s[1:]
		}
		// group digits in threes
		var b strings.Builder
		pre := len(s) % 3
		if pre == 0 {
			pre = 3
		}
		b.WriteString(s[:pre])
		for i := pre; i < len(s); i += 3 {
			b.WriteByte(' ')
			b.WriteString(s[i : i+3])
		}
		out := b.String()
		if neg {
			out = "-" + out
		}
		return out
	}

	// Unit thresholds
	abs := n
	if abs < 0 {
		abs = -abs
	}
	units := []struct {
		thresh int64
		suf    string
	}{
		{1e12, "T"},
		{1e9, "B"},
		{1e6, "M"},
		{1e3, "K"},
	}
	for _, u := range units {
		if abs >= u.thresh {
			v := float64(n) / float64(u.thresh)
			// two decimals, then trim trailing zeros and dot
			str := strconv.FormatFloat(v, 'f', 2, 64)
			str = strings.TrimRight(str, "0")
			str = strings.TrimRight(str, ".")
			return str + u.suf
		}
	}

	// smaller than 1K: just print it
	return strconv.FormatInt(n, 10)
}

// resolveFilePath resolves a filename to an absolute path
func resolveFilePath(filename string) (string, error) {
	if filepath.IsAbs(filename) {
		return filename, nil
	}

	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not get executable path: %v", err)
	}
	return filepath.Join(filepath.Dir(exePath), filename), nil
}

// LoadItemData loads data from file or fetches it via API based on file age and duration.
// If duration is 0, always fetches from API.
// If cache is false, doesn't update the file (default is true).
// Returns the parsed items as []Item.
func LoadItemData(filename string, duration int64, cache ...bool) ([]Item, error) {
	shouldCache := true
	if len(cache) > 0 {
		shouldCache = cache[0]
	}

	jsonPath, err := resolveFilePath(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not resolve file path: %v\n", err)
		return nil, err
	}

	var jsonData []byte

	// If duration is 0, always fetch
	if duration == 0 {
		jsonData, err = fetchAndCache(jsonPath, shouldCache)
		if err != nil {
			return nil, err
		}
		return loadItems(jsonData)
	}

	// Check if file exists and its modification time
	fileInfo, err := os.Stat(jsonPath)
	if os.IsNotExist(err) {
		// File doesn't exist, fetch from API
		jsonData, err = fetchAndCache(jsonPath, shouldCache)
		if err != nil {
			return nil, err
		}
		return loadItems(jsonData)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not stat file %s: %v\n", jsonPath, err)
		os.Exit(1)
	}

	// Check if file is older than duration (in seconds)
	if time.Since(fileInfo.ModTime()).Seconds() > float64(duration) {
		// File is too old, fetch from API
		jsonData, err = fetchAndCache(jsonPath, shouldCache)
		if err != nil {
			return nil, err
		}
		return loadItems(jsonData)
	}

	// File is recent enough, read from file
	jsonData, err = loadJSONFile(jsonPath)
	if err != nil {
		return nil, err
	}
	return loadItems(jsonData)
}

// fetchAndCache fetches data from API and optionally caches it
func fetchAndCache(jsonPath string, shouldCache bool) ([]byte, error) {
	payload := map[string]interface{}{"type": "itemData"}
	data := client.FetchDataOrExit(payload)
	if shouldCache {
		if err := os.WriteFile(jsonPath, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not cache data to %s: %v\n", jsonPath, err)
			return nil, err
		}
	}
	return data, nil
}

// loadJSONFile reads a JSON file and returns its contents
func loadJSONFile(filename string) ([]byte, error) {
	jsonPath, err := resolveFilePath(filename)
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", jsonPath)
	}

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %v", jsonPath, err)
	}
	return data, nil
}

// loadItems parses the jsondata and returns all items.
func loadItems(jsondata []byte) ([]Item, error) {
	if jsondata == nil {
		return nil, fmt.Errorf("jsondata is nil")
	}

	var items []Item
	if err := json.Unmarshal(jsondata, &items); err != nil {
		return nil, fmt.Errorf("failed to parse jsondata: %v", err)
	}
	return items, nil
}

// PrintJSON pretty-prints
func PrintJSON(data []byte) {
	var pretty bytes.Buffer
	json.Indent(&pretty, data, "", "  ")
	fmt.Println(pretty.String())
}

// ParseID converts arg to int
func ParseID(arg string) int {
	id, err := strconv.Atoi(arg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid ID: %s\n", arg)
		os.Exit(1)
	}
	return id
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
