package profile

import (
	"bcncli/client"
	"bcncli/common"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage user profiles",
}

func init() {
	Cmd.AddCommand(infoCmd, userCmd, inventoryCmd, statsCmd, trophiesCmd, flatinventoryCmd)

	for _, c := range []*cobra.Command{infoCmd, userCmd} {
		c.Flags().BoolP("debug", "d", false, "print raw JSON response")
		c.Flags().StringP("filter", "f", "", "comma-separated list of sections to print (e.g. farms,pets)")
		c.Flags().StringP("sort", "s", "", "sorting key for a section (e.g. farm:plant)")
	}
}

// ProfileInfo represents the detailed profile information returned by the API.
type ProfileInfo struct {
	ID                    int64             `json:"id"`
	Name                  string            `json:"name"`
	RegistrationDate      string            `json:"registrationDate"`
	Rank                  int               `json:"rank"`
	Tier                  int               `json:"tier"`
	BC                    int64             `json:"bc"`
	SP                    int64             `json:"sp"`
	KR                    int64             `json:"kr"`
	BuddyID               int64             `json:"buddyId"`
	FactionID             int64             `json:"factionId"`
	FactionTag            string            `json:"factionTag"`
	FactionJoinDate       string            `json:"factionJoinDate"`
	FactionDepositWeekly  int64             `json:"factionDepositWeekly"`
	FactionDepositTotal   int64             `json:"factionDepositTotal"`
	QuestLevel            int               `json:"questLevel"`
	QuestLevelClaimed     int               `json:"questLevelClaimed"`
	DailyClaimStreak      int               `json:"dailyClaimStreak"`
	DailyVoteStreak       int               `json:"dailyVoteStreak"`
	PetsBredDaily         int               `json:"petsBredDaily"`
	LastCaptchaDate       string            `json:"lastCaptchaDate"`
	FarmPlots             []FarmPlot        `json:"farmPlots"`
	Generators            []Generator       `json:"generators"`
	Quests                []Quest           `json:"quests"`
	Cooldowns             Cooldowns         `json:"cooldowns"`
	Effects               map[string]Effect `json:"effects"`
	Upgrades              Upgrades          `json:"upgrades"`
	EquippedFlatInventory map[string]any    `json:"equippedFlatInventory"`
	Perks                 Perks             `json:"perks"`
	PinnedItemIDs         []int64           `json:"pinnedItemIds"`
	PinnedPetIDs          []int64           `json:"pinnedPetIds"`
	AutosellLimits        map[string]int64  `json:"autosellLimits"`
	ItemReserveAmounts    map[string]int64  `json:"itemReserveAmounts"`
	Settings              Settings          `json:"settings"`
	Custom                Custom            `json:"custom"`
	DiscordServerIDs      []int64           `json:"discordServerIds"`
	BlockedBcIDs          []int64           `json:"blockedBcIds"`
	BanExpiryDate         string            `json:"banExpiryDate"`
	BanReason             *string           `json:"banReason"`
	PremiumExpiryDate     *string           `json:"premiumExpiryDate"`
	DiscordID             *string           `json:"discordId"`
	DiscordAvatarHash     *string           `json:"discordAvatarHash"`
	DiscordUsername       *string           `json:"discordUsername"`
	IsModerator           bool              `json:"isModerator"`
	Inventory             []int64           `json:"inventory"`
	Faction               Faction           `json:"faction"`
	LbPositions           LbPositions       `json:"lbPositions"`
}

// FarmPlot represents a single farm plot in the user's profile.
type FarmPlot struct {
	Level   int         `json:"level"`
	Status  PlantStatus `json:"status"`
	Boost   Boost       `json:"boost"`
	IsExtra bool        `json:"isExtra"`
}

// PlantStatus represents the planting status of a farm plot.
type PlantStatus struct {
	IsPlanted   bool  `json:"isPlanted"`
	ItemID      int64 `json:"itemId"`
	PlantedTime int64 `json:"plantedTime"`
}

// Boost represents a farm plot boost with multiplier and end time.
type Boost struct {
	Multiplier int64 `json:"multiplier"`
	EndTime    int64 `json:"endTime"`
}

// Generator represents a generator in the user's profile.
type Generator struct {
	Level   int  `json:"level"`
	IsExtra bool `json:"isExtra"`
}

// Quest represents a quest in the user's profile.
type Quest struct {
	ItemID          int64 `json:"itemId"`
	AmountRequired  int64 `json:"amountRequired"`
	AmountFulfilled int64 `json:"amountFulfilled"`
}

// Cooldowns represents the various action cooldowns in the user's profile.
type Cooldowns struct {
	Fish            int64 `json:"fish"`
	Hunt            int64 `json:"hunt"`
	Explore         int64 `json:"explore"`
	Mine            int64 `json:"mine"`
	Work            int64 `json:"work"`
	Daily           int64 `json:"daily"`
	Water           int64 `json:"water"`
	ClaimGenerators int64 `json:"claimGenerators"`
	SetBuddy        int64 `json:"setBuddy"`
	BuddyBossAttack int64 `json:"buddyBossAttack"`
	TopGgVote       int64 `json:"topGgVote"`
	Item38Use       int64 `json:"item38Use"`
}

// Effect represents a temporary effect on the user's profile.
type Effect struct {
	EndTime  int64    `json:"endTime"`
	Modifier Modifier `json:"modifier"`
}

// Modifier represents the details of an effect modifier.
type Modifier struct {
	Type       string `json:"type"`
	Action     string `json:"action,omitempty"`
	Duration   int64  `json:"duration"`
	Multiplier int64  `json:"multiplier"`
}

// Upgrades represents the user's profile upgrades.
type Upgrades struct {
	Fish            int `json:"fish"`
	FishExtra       int `json:"fishExtra"`
	Hunt            int `json:"hunt"`
	HuntExtra       int `json:"huntExtra"`
	Explore         int `json:"explore"`
	ExploreExtra    int `json:"exploreExtra"`
	Mine            int `json:"mine"`
	MineExtra       int `json:"mineExtra"`
	PetsStable      int `json:"petsStable"`
	PetsStableExtra int `json:"petsStableExtra"`
}

// Perks represents the user's profile perks.
type Perks struct {
	LowerRankCost                        int `json:"lowerRankCost"`
	LowerTierCost                        int `json:"lowerTierCost"`
	RaisePetSpace                        int `json:"raisePetSpace"`
	RaiseEquipSlots                      int `json:"raiseEquipSlots"`
	RaisePetMaxTier                      int `json:"raisePetMaxTier"`
	LowerPetBreedCost                    int `json:"lowerPetBreedCost"`
	RaiseCoinflipLimit                   int `json:"raiseCoinflipLimit"`
	RaiseWorkBonusChance                 int `json:"raiseWorkBonusChance"`
	RaiseFarmCropsDieTime                int `json:"raiseFarmCropsDieTime"`
	LowerWaterFarmCooldown               int `json:"lowerWaterFarmCooldown"`
	RaiseGeneratorIdleTime               int `json:"raiseGeneratorIdleTime"`
	RaisePetEnergyCapacity               int `json:"raisePetEnergyCapacity"`
	RaiseMaxSameItemPlanted              int `json:"raiseMaxSameItemPlanted"`
	RaiseRareItemMultiplier              int `json:"raiseRareItemMultiplier"`
	RaiseFarmWaterByproducts             int `json:"raiseFarmWaterByproducts"`
	RaiseToolAugmentationSlots           int `json:"raiseToolAugmentationSlots"`
	RaisePetCravingXpMultiplier          int `json:"raisePetCravingXpMultiplier"`
	RaisePetFeedAdditionalItemOutput     int `json:"raisePetFeedAdditionalItemOutput"`
	RaiseChanceToIgnoreCooldownForAction int `json:"raiseChanceToIgnoreCooldownForAction"`
	RaiseFarmHarvestAdditionalItemOutput int `json:"raiseFarmHarvestAdditionalItemOutput"`
}

// Settings represents the user's profile settings.
type Settings struct {
	ProfileShowStatID     *int    `json:"profileShowStatId"`
	Title                 *string `json:"title"`
	SyncDiscordName       bool    `json:"syncDiscordName"`
	PublicDiscordProfile  bool    `json:"publicDiscordProfile"`
	DiscordPingOnResponse bool    `json:"discordPingOnResponse"`
}

// Custom represents the user's profile customizations.
type Custom struct {
	ProfileHideAvatar          bool    `json:"profileHideAvatar"`
	ProfileHideTitleName       bool    `json:"profileHideTitleName"`
	ProfileUseChatEmblemEmoji  bool    `json:"profileUseChatEmblemEmoji"`
	ProfileBackground          *string `json:"profileBackground"`
	ChatEmblemEmoji            *string `json:"chatEmblemEmoji"`
	ChatUsernameColor          *string `json:"chatUsernameColor"`
	ChatUsernameStyle          *string `json:"chatUsernameStyle"`
	ChatMessageBackgroundColor *string `json:"chatMessageBackgroundColor"`
}

// BoostStep represents a step in the faction boost system.
type BoostStep struct {
	LastChange int64 `json:"lastChange"`
	Amount     int   `json:"amount"`
}

// CustomizationSettings represents the customization options for a faction.
type CustomizationSettings struct {
	EmblemEmoji string  `json:"emblemEmoji"`
	TagColor    string  `json:"tagColor"`
	NameColor   *string `json:"nameColor"`
	NameStyle   string  `json:"nameStyle"`
}

// Faction represents a faction in the game.
type Faction struct {
	ID                     int64                 `json:"id"`
	Tag                    string                `json:"tag"`
	Name                   string                `json:"name"`
	OwnerBcID              int64                 `json:"ownerBcId"`
	RankOverrides          map[string]int        `json:"rankOverrides"`
	IsRecruiting           bool                  `json:"isRecruiting"`
	About                  string                `json:"about"`
	Motd                   string                `json:"motd"`
	UnsyncedFp             int64                 `json:"unsyncedFp"`
	LastFpSync             string                `json:"lastFpSync"`
	BoostSteps             map[string]BoostStep  `json:"boostSteps"`
	Halls                  int64                 `json:"halls"`
	FpDepositedMonthly     int64                 `json:"fpDepositedMonthly"`
	FpDepositedTotal       int64                 `json:"fpDepositedTotal"`
	CustomizationSettings  CustomizationSettings `json:"customizationSettings"`
	OwnerPremiumExpiryDate string                `json:"ownerPremiumExpiryDate"`
	MemberCount            int                   `json:"memberCount"`
	PendingRequests        int                   `json:"pendingRequests"`
}

// LbPositions represents the leaderboard positions for a profile.
type LbPositions struct {
	Rank                   int `json:"rank"`
	IncomeDaily            int `json:"incomeDaily"`
	NetCoinflipProfitDaily int `json:"netCoinflipProfitDaily"`
}

var infoCmd = &cobra.Command{
	Use:   "info [id]",
	Short: "Fetch profile info (detailed)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		executeProfileCmd(cmd, args, "profile")
	},
}

// renderProfileInfoTable prints a short summary using the standard
// text/tabwriter for tidy alignment.
func renderProfileInfoTable(p ProfileInfo) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	defer tw.Flush()

	row := func(label, value string) {
		fmt.Fprintf(tw, "%s:	%s\n", label, value)
	}

	row("ID", strconv.FormatInt(p.ID, 10))
	row("Name", p.Name)
	row("Registered", p.RegistrationDate)
	row("Rank", strconv.Itoa(p.Rank))
	row("Tier", strconv.Itoa(p.Tier))
	row("BC", fmt.Sprintf("%d", p.BC))
	row("SP", fmt.Sprintf("%d", p.SP))
	row("KR", fmt.Sprintf("%d", p.KR))
	row("Faction", fmt.Sprintf("%s (ID %d)", p.FactionTag, p.FactionID))
	row("Quest Level", strconv.Itoa(p.QuestLevel))
	row("Daily Streak", strconv.Itoa(p.DailyClaimStreak))
}

var userCmd = &cobra.Command{
	Use:   "user [id]",
	Short: "Fetch user details (alias for info)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		executeProfileCmd(cmd, args, "user")
	},
}

// executeProfileCmd is shared by infoCmd and userCmd to avoid duplication.
func executeProfileCmd(cmd *cobra.Command, args []string, payloadType string) {
	userID := common.ParseID(args[0])

	// handle debug flag early so we do not unmarshal twice
	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		payload := map[string]any{"type": payloadType, "id": userID}
		raw := client.FetchDataOrExit(payload)
		common.PrintJSON(raw)
		return
	}

	// fetch + unmarshal
	payload := map[string]any{"type": payloadType, "id": userID}
	raw := client.FetchDataOrExit(payload)
	var profile ProfileInfo
	if err := json.Unmarshal(raw, &profile); err != nil {
		fmt.Println("error decoding profile info:", err)
		os.Exit(1)
	}

	// parse flags
	filterFlag, _ := cmd.Flags().GetString("filter")
	sortFlag, _ := cmd.Flags().GetString("sort")
	filters := parseFilter(filterFlag)

	renderProfile(profile, filters, sortFlag)
}

// ===============================
// Rendering helpers
// ===============================

type sectionWriter struct {
	tw *tabwriter.Writer
}

func newSectionWriter() *sectionWriter {
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	return &sectionWriter{tw: tw}
}

func (sw *sectionWriter) flush() { _ = sw.tw.Flush() }

func (sw *sectionWriter) title(t string) { fmt.Fprintf(sw.tw, "\n=== %s ===\n", strings.ToUpper(t)) }

func (sw *sectionWriter) row(k, v string) { fmt.Fprintf(sw.tw, "%s:\t%s\n", k, v) }

// renderProfile prints every possible field of ProfileInfo.
// If filters is non-empty only the requested sections are rendered.
func renderProfile(p ProfileInfo, filters map[string]bool, sortFlag string) {
	sw := newSectionWriter()
	defer sw.flush()

	want := func(name string) bool {
		return len(filters) == 0 || filters[strings.ToLower(name)]
	}

	// Basic section
	if want("basic") {
		sw.title("Basic")
		sw.row("ID", strconv.FormatInt(p.ID, 10))
		sw.row("Name", p.Name)
		sw.row("Registered", p.RegistrationDate)
		sw.row("Rank", strconv.Itoa(p.Rank))
		sw.row("Tier", strconv.Itoa(p.Tier))
		sw.row("BC", fmt.Sprintf("%d", p.BC))
		sw.row("SP", fmt.Sprintf("%d", p.SP))
		sw.row("KR", fmt.Sprintf("%d", p.KR))
		sw.row("Quest Level", strconv.Itoa(p.QuestLevel))
		sw.row("Daily Streak", strconv.Itoa(p.DailyClaimStreak))
	}

	// Faction
	if want("faction") && p.FactionID != 0 {
		sw.title("Faction")
		sw.row("Tag", p.Faction.Tag)
		sw.row("Name", p.Faction.Name)
		sw.row("Member Count", strconv.Itoa(p.Faction.MemberCount))
		sw.row("Owner", strconv.FormatInt(p.Faction.OwnerBcID, 10))
		sw.row("About", p.Faction.About)
		sw.row("MOTD", p.Faction.Motd)
	}

	// Farm Plots
	if want("farms") {
		sw.title("Farm Plots")
		plots := append([]FarmPlot(nil), p.FarmPlots...) // copy to avoid mutation
		sortFarmPlots(plots, sortFlag)
		for i, fp := range plots {
			prefix := fmt.Sprintf("Plot %d", i+1)
			sw.row(prefix+" Level", strconv.Itoa(fp.Level))
			sw.row(prefix+" Extra", strconv.FormatBool(fp.IsExtra))
			sw.row(prefix+" Planted", strconv.FormatBool(fp.Status.IsPlanted))
			if fp.Status.IsPlanted {
				sw.row(prefix+" ItemID", strconv.FormatInt(fp.Status.ItemID, 10))
				sw.row(prefix+" PlantedTime", strconv.FormatInt(fp.Status.PlantedTime, 10))
			}
			if fp.Boost.Multiplier > 0 {
				sw.row(prefix+" Boost x", fmt.Sprintf("%d", fp.Boost.Multiplier))
				sw.row(prefix+" BoostEnd", strconv.FormatInt(fp.Boost.EndTime, 10))
			}
		}
	}

	// Generators
	if want("generators") {
		sw.title("Generators")
		for i, g := range p.Generators {
			prefix := fmt.Sprintf("Gen %d", i+1)
			sw.row(prefix+" Level", strconv.Itoa(g.Level))
			sw.row(prefix+" Extra", strconv.FormatBool(g.IsExtra))
		}
	}

	// Quests
	if want("quests") {
		sw.title("Quests")
		for i, q := range p.Quests {
			prefix := fmt.Sprintf("Quest %d", i+1)
			sw.row(prefix+" ItemID", strconv.FormatInt(q.ItemID, 10))
			sw.row(prefix+" Required", strconv.FormatInt(q.AmountRequired, 10))
			sw.row(prefix+" Fulfilled", strconv.FormatInt(q.AmountFulfilled, 10))
		}
	}

	// Cooldowns
	if want("cooldowns") {
		sw.title("Cooldowns (Unix ms)")
		sw.row("Fish", strconv.FormatInt(p.Cooldowns.Fish, 10))
		sw.row("Hunt", strconv.FormatInt(p.Cooldowns.Hunt, 10))
		sw.row("Explore", strconv.FormatInt(p.Cooldowns.Explore, 10))
		sw.row("Mine", strconv.FormatInt(p.Cooldowns.Mine, 10))
		sw.row("Work", strconv.FormatInt(p.Cooldowns.Work, 10))
		sw.row("Daily", strconv.FormatInt(p.Cooldowns.Daily, 10))
		sw.row("Water", strconv.FormatInt(p.Cooldowns.Water, 10))
		sw.row("ClaimGenerators", strconv.FormatInt(p.Cooldowns.ClaimGenerators, 10))
	}

	// Effects / Modifiers
	if want("effects") {
		sw.title("Effects")
		keys := make([]string, 0, len(p.Effects))
		for k := range p.Effects {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			e := p.Effects[k]
			sw.row(k+" End", strconv.FormatInt(e.EndTime, 10))
			sw.row(k+" Type", e.Modifier.Type)
			if e.Modifier.Action != "" {
				sw.row(k+" Action", e.Modifier.Action)
			}
			sw.row(k+" Mult x", strconv.FormatInt(e.Modifier.Multiplier, 10))
		}
	}

	// Upgrades
	if want("upgrades") {
		sw.title("Upgrades")
		sw.row("Fish", strconv.Itoa(p.Upgrades.Fish))
		sw.row("FishExtra", strconv.Itoa(p.Upgrades.FishExtra))
		sw.row("Hunt", strconv.Itoa(p.Upgrades.Hunt))
		sw.row("HuntExtra", strconv.Itoa(p.Upgrades.HuntExtra))
		sw.row("Explore", strconv.Itoa(p.Upgrades.Explore))
		sw.row("ExploreExtra", strconv.Itoa(p.Upgrades.ExploreExtra))
		sw.row("Mine", strconv.Itoa(p.Upgrades.Mine))
		sw.row("MineExtra", strconv.Itoa(p.Upgrades.MineExtra))
		sw.row("PetsStable", strconv.Itoa(p.Upgrades.PetsStable))
		sw.row("PetsStableExtra", strconv.Itoa(p.Upgrades.PetsStableExtra))
	}

	// Perks
	if want("perks") {
		sw.title("Perks")
		sw.row("LowerRankCost", strconv.Itoa(p.Perks.LowerRankCost))
		sw.row("LowerTierCost", strconv.Itoa(p.Perks.LowerTierCost))
		sw.row("RaisePetSpace", strconv.Itoa(p.Perks.RaisePetSpace))
		sw.row("RaiseEquipSlots", strconv.Itoa(p.Perks.RaiseEquipSlots))
		// (add remaining perks as needed)
	}

	// Settings
	if want("settings") {
		sw.title("Settings")
		if p.Settings.ProfileShowStatID != nil {
			sw.row("ProfileShowStatID", strconv.Itoa(*p.Settings.ProfileShowStatID))
		}
		sw.row("SyncDiscordName", fmt.Sprintf("%t", p.Settings.SyncDiscordName))
		sw.row("PublicDiscordProfile", fmt.Sprintf("%t", p.Settings.PublicDiscordProfile))
		sw.row("DiscordPingOnResponse", fmt.Sprintf("%t", p.Settings.DiscordPingOnResponse))
	}

	// Custom
	if want("custom") {
		sw.title("Custom")
		sw.row("HideAvatar", fmt.Sprintf("%t", p.Custom.ProfileHideAvatar))
		sw.row("HideTitleName", fmt.Sprintf("%t", p.Custom.ProfileHideTitleName))
		sw.row("UseChatEmblemEmoji", fmt.Sprintf("%t", p.Custom.ProfileUseChatEmblemEmoji))
		if p.Custom.ProfileBackground != nil {
			sw.row("Background", *p.Custom.ProfileBackground)
		}
	}
}

// sortFarmPlots sorts plots in-place when the user provided a sort flag like "farm:plant".
// Currently supported keys: plant, level.
func sortFarmPlots(plots []FarmPlot, sortFlag string) {
	if len(plots) == 0 || sortFlag == "" {
		return
	}
	parts := strings.Split(sortFlag, ":")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "farm" {
		return // not meant for farms
	}

	key := strings.ToLower(parts[1])
	switch key {
	case "plant", "item":
		sort.SliceStable(plots, func(i, j int) bool {
			return plots[i].Status.ItemID < plots[j].Status.ItemID
		})
	case "level":
		sort.SliceStable(plots, func(i, j int) bool {
			return plots[i].Level < plots[j].Level
		})
	}
}

// ===============================
// Helper utilities
// ===============================

// parseFilter converts comma-separated filter string to a set.
func parseFilter(raw string) map[string]bool {
	if raw == "" {
		return nil
	}
	set := make(map[string]bool)
	for _, f := range strings.Split(raw, ",") {
		f = strings.ToLower(strings.TrimSpace(f))
		if f == "" {
			continue
		}
		set[f] = true
	}
	return set
}

var inventoryCmd = &cobra.Command{
	Use:   "inventory [id]",
	Short: "Fetch inventory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := common.ParseID(args[0])
		payload := map[string]interface{}{"type": "inventory", "id": id}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}

var flatinventoryCmd = &cobra.Command{
	Use:   "flatinventory [id]",
	Short: "Fetch flat inventory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := common.ParseID(args[0])
		payload := map[string]interface{}{"type": "flatInventory", "id": id}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats [id]",
	Short: "Fetch stats",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := common.ParseID(args[0])
		payload := map[string]interface{}{"type": "stats", "id": id}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}

var trophiesCmd = &cobra.Command{
	Use:   "trophies [id]",
	Short: "Fetch trophies",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := common.ParseID(args[0])
		payload := map[string]interface{}{"type": "trophies", "id": id}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}
