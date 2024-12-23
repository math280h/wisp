package shared

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	BotToken       = flag.String("token", "", "Bot access token")
	GuildID        = flag.String("guild", "", "Guild ID")
	ReportCategory = flag.String("category", "", "Category ID")
	ArchiveChannel = flag.String("archive", "", "Archive channel ID")
	LogChannel     = flag.String("log", "", "Log channel ID")
	AlertChannel   = flag.String("alert", "", "Alert channel ID")
	PrettyLogs     = flag.Bool("pretty", false, "Pretty logs")

	// Moderation Settings
	WarnPoints   = flag.Int("warnpoints", 10, "Number of points to warn a user")
	StrikePoints = flag.Int("strikepoints", 50, "Number of points to strike a user")
	MaxPoints    = flag.Int("maxpoints", 100, "Number of points to ban a user")

	// Suggestions
	SuggestionChannel = flag.String("suggestion", "", "Suggestion channel ID")
)

var GuildName string = "Unknown"

func Init() {
	flag.Parse()

	env_err := godotenv.Load()
	if env_err != nil {
		log.Fatal("Error loading .env file")
	}

	// If BotToken is not provided, use the one from the .env file
	if *BotToken == "" {
		*BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	}
	// If GuildID is not provided, use the one from the .env file
	if *GuildID == "" {
		*GuildID = os.Getenv("DISCORD_GUILD_ID")
	}
	// If ReportCategory is not provided, use the one from the .env file
	if *ReportCategory == "" {
		*ReportCategory = os.Getenv("DISCORD_GUILD_REPORT_CATEGORY_ID")
	}
	// If ArchiveChannel is not provided, use the one from the .env file
	if *ArchiveChannel == "" {
		*ArchiveChannel = os.Getenv("DISCORD_GUILD_ARCHIVE_CHANNEL")
	}
	// If LogChannel is not provided, use the one from the .env file
	if *LogChannel == "" {
		*LogChannel = os.Getenv("DISCORD_GUILD_LOG_CHANNEL")
	}
	// If AlertChannel is not provided, use the one from the .env file
	if *AlertChannel == "" {
		*AlertChannel = os.Getenv("DISCORD_GUILD_ALERT_CHANNEL")
	}
	// If PrettyLogs is not provided, use the one from the .env file
	if !*PrettyLogs {
		pretty_logs := os.Getenv("PRETTY_LOGS")
		if pretty_logs == "true" {
			*PrettyLogs = true
		}
	}

	// Moderation Settings

	// If WarnPoints is not provided, use the one from the .env file
	if *WarnPoints == 10 {
		warn_points := os.Getenv("WARN_POINTS")
		if warn_points != "" {
			*WarnPoints, _ = strconv.Atoi(warn_points)
		}
	}
	// If StrikePoints is not provided, use the one from the .env file
	if *StrikePoints == 50 {
		strike_points := os.Getenv("STRIKE_POINTS")
		if strike_points != "" {
			*StrikePoints, _ = strconv.Atoi(strike_points)
		}
	}
	// If MaxPoints is not provided, use the one from the .env file
	if *MaxPoints == 100 {
		max_points := os.Getenv("MAX_POINTS")
		if max_points != "" {
			*MaxPoints, _ = strconv.Atoi(max_points)
		}
	}

	// Suggestions

	// If SuggestionChannel is not provided, use the one from the .env file
	if *SuggestionChannel == "" {
		*SuggestionChannel = os.Getenv("DISCORD_SUGGESTION_CHANNEL")
	}
}

func SetGuildName(s *discordgo.Session) {
	GuildName = s.State.Guilds[0].Name
}
