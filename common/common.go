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

// FoodItem represents an entry of food item, with name and energy value.
type FoodItem struct {
	Name   string
	Energy int
}

// AllFoodItems is a list of all food items with their energy values.
var AllFoodItems = []FoodItem{
	{"Seaweed", 25},
	{"Sardine", 50},
	{"Exotic Bean", 150},
	{"Prawn", 150},
	{"Red Mushroom", 200},
	{"Bird Nest", 300},
	{"Jellyfish", 350},
	{"Soybean", 500},
	{"Milk", 500},
	{"Prime Steak", 600},
	{"Ocean Crab", 750},
	{"Blueberry", 750},
	{"Golden Wheat", 1000},
	{"Russet Potato", 1000},
	{"Blowfish", 2500},
	{"Electric Eel", 5000},
	{"Strawberry", 7500},
	{"Kiwi", 12500},
	{"Seafood Salad", 23000},
	{"Great White", 25000},
	{"Mango", 25000},
	{"Hearty Burger", 32500},
	{"Melon", 50000},
	{"Warm Broth", 58100},
	{"Pearled Oyster", 100000},
	{"Stone Soup", 268100},
	{"Coconut", 750000},
	{"Giant Squid", 1000000},
	{"Pumpkin", 7500000},
}

// PetBoostItem represents a pet boost with its BC worth and effect.
type PetBoostItem struct {
	Name   string // Human-readable name of the boost
	Worth  int    // Value in BC, stored as a plain integer
	Effect string // Description of the boost effect
}

// AllPetBoostItems holds our collection of pet boost items.
var AllPetBoostItems = []PetBoostItem{
	{"Fragrant Dogrose", 2500000, "2× Pet Adventure Speed (2h)"},
	{"Mystical Rowan", 25000000, "4× Pet Adventure Speed (2h)"},
	{"Legendary Aguaje", 250000000, "8× Pet Adventure Speed (2h)"},
	{"Magic Token", 1000000000, "10× Pet Adventure Speed (1d)"},
}

// PetData represents a creature with an icon, name, and category.
type PetData struct {
	Icon     string // Emoji or icon representation
	Name     string // Human-readable name, without the icon
	Category string // One of "Fish", "Hunt", "Explore", "Mine"
}

// AllPetTypes is the full list of available pets.
var AllPetTypes = []PetData{
	// Fish
	{"🐬", "Dolphin", "Fish"},
	{"🦦", "Otter", "Fish"},
	{"🪿", "Goose", "Fish"},
	{"🦭", "Seal", "Fish"},
	{"🐳", "Whale", "Fish"},
	{"🐢", "Turtle", "Fish"},
	{"", "Dragon", "Fish"},

	// Hunt
	{"🦅", "Eagle", "Hunt"},
	{"🐅", "Tiger", "Hunt"},
	{"🦍", "Gorilla", "Hunt"},
	{"🐊", "Crocodile", "Hunt"},
	{"🐍", "Snake", "Hunt"},
	{"", "Scorpion", "Hunt"},
	{"", "Phoenix", "Hunt"},

	// Explore
	{"🐩", "Poodle", "Explore"},
	{"🐕", "Dog", "Explore"},
	{"🐎", "Mustang", "Explore"},
	{"🐖", "Pig", "Explore"},
	{"🦚", "Peacock", "Explore"},
	{"🫏", "Donkey", "Explore"},
	{"🐂", "Ox", "Explore"},
	{"🐓", "Junglefowl", "Explore"},
	{"🐇", "Rabbit", "Explore"},
	{"🕊️", "Dove", "Explore"},
	{"🦘", "Kangaroo", "Explore"},
	{"", "Visitor", "Explore"},

	// Mine
	{"🦇", "Bat", "Mine"},
	{"🐀", "Rat", "Mine"},
	{"🐌", "Snail", "Mine"},
	{"🦎", "Lizard", "Mine"},
	{"", "Invader", "Mine"},
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

// TimeUntilISO8601 takes an RFC3339 timestamp (EpochToISO8601),
// and returns the time from now until that instant in a human‑readable form.
// If the time is now or in the past, it returns "0".
func TimeUntilISO8601(iso string) string {
	// parse the incoming timestamp
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return "-" // or handle parse error as you prefer
	}
	now := time.Now().UTC()
	diff := t.Sub(now)

	// if zero or negative, we’re done
	if diff <= 0 {
		return "0"
	}

	// break diff down into components
	totalSeconds := int64(diff.Seconds())
	weeks := totalSeconds / (7 * 24 * 3600)
	totalSeconds %= 7 * 24 * 3600
	days := totalSeconds / (24 * 3600)
	totalSeconds %= 24 * 3600
	hours := totalSeconds / 3600
	totalSeconds %= 3600
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	// build the human‑readable parts
	var parts []string
	if weeks > 0 {
		parts = append(parts, fmt.Sprintf("%dw", weeks))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	// always show seconds if it’s the only unit, otherwise only if >0
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	return strings.Join(parts, " ")
}

// ElapsedSinceISO8601 takes an RFC3339 timestamp (EpochToISO8601),
// and returns the time elapsed from that instant until now in a human‑readable form.
// If the time is in the future or parsing fails, it returns "0".
func ElapsedSinceISO8601(iso string) string {
	// parse the incoming timestamp
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return "0" // or handle error otherwise
	}
	now := time.Now().UTC()
	diff := now.Sub(t)

	// if zero or negative (i.e. t is in the future), we’re done
	if diff <= 0 {
		return "0"
	}

	// break diff down into components
	totalSeconds := int64(diff.Seconds())
	weeks := totalSeconds / (7 * 24 * 3600)
	totalSeconds %= 7 * 24 * 3600
	days := totalSeconds / (24 * 3600)
	totalSeconds %= 24 * 3600
	hours := totalSeconds / 3600
	totalSeconds %= 3600
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	// build the human‑readable parts
	var parts []string
	if weeks > 0 {
		parts = append(parts, fmt.Sprintf("%dw", weeks))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	// always show seconds if it’s the only unit, otherwise only if >0
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	return strings.Join(parts, " ")
}

// LookUpItemName finds the item name by ID in the provided items slice.
func LookUpItemName(id int, items []Item) string {
	for _, item := range items {
		if item.ID == id {
			return item.Name
		}
	}
	return fmt.Sprintf("Unknown Item ID %d", id)
}

// GetEnergy returns the energy value for the given item name.
// If the item is not found, it returns 0.
func GetEnergy(name string) int {
	for _, it := range AllFoodItems {
		if strings.EqualFold(it.Name, name) {
			return it.Energy
		}
	}
	return 0
}

// GetPetBoostDetails looks up a boost by name (case-insensitive).
// It returns the worth (in BC) and effect string.
// If the boost isn't found, it returns 0 and an empty string.
func GetPetBoostDetails(name string) (int, string) {
	for _, boost := range AllPetBoostItems {
		if strings.EqualFold(boost.Name, name) {
			return boost.Worth, boost.Effect
		}
	}
	return 0, ""
}

// GetPetCategory returns the category of the pet with the given name (case-insensitive).
// If the pet is not found, it returns an empty string.
func GetPetCategory(name string) string {
	for _, pet := range AllPetTypes {
		if strings.EqualFold(pet.Name, name) {
			return pet.Category
		}
	}
	return ""
}

// GetPetsByCategory returns a slice of all pets in the given category (case-insensitive).
// If no pets are found, it returns an empty slice.
func GetPetsByCategory(category string) []PetData {
	var results []PetData
	for _, pet := range AllPetTypes {
		if strings.EqualFold(pet.Category, category) {
			results = append(results, pet)
		}
	}
	return results
}
