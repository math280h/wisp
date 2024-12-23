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
	BotToken       = flag.String("token", "", "Bot access token")     //nolint:gochecknoglobals,lll // This is a flag shared across the application
	GuildID        = flag.String("guild", "", "Guild ID")             //nolint:gochecknoglobals,lll // This is a flag shared across the application
	ReportCategory = flag.String("category", "", "Category ID")       //nolint:gochecknoglobals,lll // This is a flag shared across the application
	ArchiveChannel = flag.String("archive", "", "Archive channel ID") //nolint:gochecknoglobals,lll // This is a flag shared across the application
	LogChannel     = flag.String("log", "", "Log channel ID")         //nolint:gochecknoglobals,lll // This is a flag shared across the application
	AlertChannel   = flag.String("alert", "", "Alert channel ID")     //nolint:gochecknoglobals,lll // This is a flag shared across the application
	MutedRole      = flag.String("muted", "", "Muted role ID")        //nolint:gochecknoglobals,lll // This is a flag shared across the application
	PrettyLogs     = flag.Bool("pretty", false, "Pretty logs")        //nolint:gochecknoglobals,lll // This is a flag shared across the application

	// Moderation Settings.
	WarnPoints   = flag.Int("warnpoints", 10, "Number of points to warn a user")     //nolint:gochecknoglobals,lll // This is a flag shared across the application
	StrikePoints = flag.Int("strikepoints", 50, "Number of points to strike a user") //nolint:gochecknoglobals,lll // This is a flag shared across the application
	MaxPoints    = flag.Int("maxpoints", 100, "Number of points to ban a user")      //nolint:gochecknoglobals,lll // This is a flag shared across the application

	// Suggestions.
	SuggestionChannel = flag.String("suggestion", "", "Suggestion channel ID") //nolint:gochecknoglobals,lll // This is a flag shared across the application
)

var GuildName = "Unknown" //nolint:gochecknoglobals // This is the name of the guild

func Init() { //nolint:gocognit // This function is responsible for initializing the shared flags
	flag.Parse()

	envErr := godotenv.Load()
	if envErr != nil {
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
	// If MutedRole is not provided, use the one from the .env file
	if *MutedRole == "" {
		*MutedRole = os.Getenv("DISCORD_GUILD_MUTED_ROLE")
	}
	// If PrettyLogs is not provided, use the one from the .env file
	if !*PrettyLogs {
		envPrettyLogs := os.Getenv("PRETTY_LOGS")
		if envPrettyLogs == "true" {
			*PrettyLogs = true
		}
	}

	// Moderation Settings

	// If WarnPoints is not provided, use the one from the .env file
	if *WarnPoints == 10 {
		envWarnPoints := os.Getenv("WARN_POINTS")
		if envWarnPoints != "" {
			*WarnPoints, _ = strconv.Atoi(envWarnPoints)
		}
	}
	// If StrikePoints is not provided, use the one from the .env file
	if *StrikePoints == 50 {
		envStrikePoints := os.Getenv("STRIKE_POINTS")
		if envStrikePoints != "" {
			*StrikePoints, _ = strconv.Atoi(envStrikePoints)
		}
	}
	// If MaxPoints is not provided, use the one from the .env file
	if *MaxPoints == 100 {
		envMaxPoints := os.Getenv("MAX_POINTS")
		if envMaxPoints != "" {
			*MaxPoints, _ = strconv.Atoi(envMaxPoints)
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
