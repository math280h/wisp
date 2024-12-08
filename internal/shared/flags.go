package shared

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	BotToken       = flag.String("token", "", "Bot access token")
	GuildID        = flag.String("guild", "", "Guild ID")
	ReportCategory = flag.String("category", "", "Category ID")
	ArchiveChannel = flag.String("archive", "", "Archive channel ID")
	LogChannel     = flag.String("log", "", "Log channel ID")
)

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
}
