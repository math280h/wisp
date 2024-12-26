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
	// General.
	BotToken = flag.String("token", "", "Bot access token") //nolint:gochecknoglobals,lll // This is a flag shared across the application

	// Guild Information.
	GuildID = flag.String("guild", "", "Guild ID") //nolint:gochecknoglobals,lll // This is a flag shared across the application
	// Channels.
	ReportCategory    = flag.String("category", "", "Category ID")             //nolint:gochecknoglobals,lll // This is a flag shared across the application
	ArchiveChannel    = flag.String("archive", "", "Archive channel ID")       //nolint:gochecknoglobals,lll // This is a flag shared across the application
	LogChannel        = flag.String("log", "", "Log channel ID")               //nolint:gochecknoglobals,lll // This is a flag shared across the application
	AlertChannel      = flag.String("alert", "", "Alert channel ID")           //nolint:gochecknoglobals,lll // This is a flag shared across the application
	HistoryChannel    = flag.String("history", "", "History channel ID")       //nolint:gochecknoglobals,lll // This is a flag shared across the application
	SuggestionChannel = flag.String("suggestion", "", "Suggestion channel ID") //nolint:gochecknoglobals,lll // This is a flag shared across the application
	// Roles.
	MutedRole = flag.String("muted", "", "Muted role ID") //nolint:gochecknoglobals,lll // This is a flag shared across the application

	// Moderation Settings.
	WarnPoints   = flag.Int("warnpoints", 10, "Number of points to warn a user")     //nolint:gochecknoglobals,lll // This is a flag shared across the application
	StrikePoints = flag.Int("strikepoints", 50, "Number of points to strike a user") //nolint:gochecknoglobals,lll // This is a flag shared across the application
	MaxPoints    = flag.Int("maxpoints", 100, "Number of points to ban a user")      //nolint:gochecknoglobals,lll // This is a flag shared across the application

	// Feature Flags.
	MessageHistoryEnabled  = flag.Bool("history_enabled", false, "Enable message history")   //nolint:gochecknoglobals,lll // This is a flag shared across the application
	NicknameHistoryEnabled = flag.Bool("nickname_enabled", false, "Enable nickname history") //nolint:gochecknoglobals,lll // This is a flag shared across the application

	// Bot Information.
	PrettyLogs = flag.Bool("pretty", false, "Pretty logs") //nolint:gochecknoglobals,lll // This is a flag shared across the application
)

var GuildName = "Unknown" //nolint:gochecknoglobals // This is the name of the guild

const enabled = "true"

func Init() { //nolint:gocognit // This function is responsible for initializing the shared flags
	flag.Parse()

	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	// General
	if *BotToken == "" {
		*BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	}

	// Guild Information
	if *GuildID == "" {
		*GuildID = os.Getenv("DISCORD_GUILD_ID")
	}
	// Channels
	if *ReportCategory == "" {
		*ReportCategory = os.Getenv("DISCORD_GUILD_REPORT_CATEGORY_ID")
	}
	if *ArchiveChannel == "" {
		*ArchiveChannel = os.Getenv("DISCORD_GUILD_ARCHIVE_CHANNEL")
	}
	if *LogChannel == "" {
		*LogChannel = os.Getenv("DISCORD_GUILD_LOG_CHANNEL")
	}
	if *AlertChannel == "" {
		*AlertChannel = os.Getenv("DISCORD_GUILD_ALERT_CHANNEL")
	}
	if *HistoryChannel == "" {
		*HistoryChannel = os.Getenv("DISCORD_GUILD_HISTORY_CHANNEL")
	}
	if *SuggestionChannel == "" {
		*SuggestionChannel = os.Getenv("DISCORD_SUGGESTION_CHANNEL")
	}
	// Roles
	if *MutedRole == "" {
		*MutedRole = os.Getenv("DISCORD_GUILD_MUTED_ROLE")
	}

	// Moderation Settings
	if *WarnPoints == 10 {
		envWarnPoints := os.Getenv("WARN_POINTS")
		if envWarnPoints != "" {
			*WarnPoints, _ = strconv.Atoi(envWarnPoints)
		}
	}
	if *StrikePoints == 50 {
		envStrikePoints := os.Getenv("STRIKE_POINTS")
		if envStrikePoints != "" {
			*StrikePoints, _ = strconv.Atoi(envStrikePoints)
		}
	}
	if *MaxPoints == 100 {
		envMaxPoints := os.Getenv("MAX_POINTS")
		if envMaxPoints != "" {
			*MaxPoints, _ = strconv.Atoi(envMaxPoints)
		}
	}

	// Feature Flags
	if !*MessageHistoryEnabled {
		envMessageHistoryEnabled := os.Getenv("MESSAGE_HISTORY_ENABLED")
		if envMessageHistoryEnabled == enabled {
			*MessageHistoryEnabled = true
		}
	}
	if !*NicknameHistoryEnabled {
		envNicknameHistoryEnabled := os.Getenv("NICKNAME_HISTORY_ENABLED")
		if envNicknameHistoryEnabled == enabled {
			*NicknameHistoryEnabled = true
		}
	}

	// Bot Information
	if !*PrettyLogs {
		envPrettyLogs := os.Getenv("PRETTY_LOGS")
		if envPrettyLogs == enabled {
			*PrettyLogs = true
		}
	}
}

func SetGuildName(s *discordgo.Session) {
	GuildName = s.State.Guilds[0].Name
}
