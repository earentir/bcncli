package profile

import (
	"bcncli/client"
	"bcncli/common"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage user profiles",
}

func init() {
	Cmd.AddCommand(infoCmd, userCmd, inventoryCmd, statsCmd, trophiesCmd, flatinventoryCmd)
	infoCmd.Flags().BoolP("debug", "d", false, "print raw JSON response")
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
	Short: "Fetch profile info",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userID := common.ParseID(args[0])

		payload := map[string]any{"type": "profile", "id": userID}
		raw := client.FetchDataOrExit(payload)

		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			common.PrintJSON(raw)
			return
		}

		var responce ProfileInfo
		if err := json.Unmarshal(raw, &responce); err != nil {
			fmt.Println("Error decoding profile info:", err)
			os.Exit(1)
		}

		renderProfileInfoTable(responce)
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
	Short: "Fetch user details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userID := common.ParseID(args[0])
		payload := map[string]any{"type": "user", "id": userID}
		raw := client.FetchDataOrExit(payload)

		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			common.PrintJSON(raw)
			return
		}

		var p ProfileInfo
		if err := json.Unmarshal(raw, &p); err != nil {
			fmt.Println("Error decoding user details:", err)
			os.Exit(1)
		}

		renderProfileInfoTable(p)
	},
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
