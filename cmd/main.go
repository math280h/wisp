package main

import (
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"math280h/wisp/internal/db"
	"math280h/wisp/internal/history"
	"math280h/wisp/internal/moderation"
	"math280h/wisp/internal/reports"
	"math280h/wisp/internal/shared"
	"math280h/wisp/internal/suggestions"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

var s *discordgo.Session //nolint:gochecknoglobals // This is the Discord session

var (
	commands = []*discordgo.ApplicationCommand{} //nolint:gochecknoglobals // This is a list of commands for Discord

	commandHandlers = map[string]func( //nolint:gochecknoglobals // This is a map of commands to their handlers
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	){}
)

func main() {
	shared.Init()

	// Add moderation commands
	commands = append(commands, moderation.Commands...)
	for k, v := range moderation.CommandHandlers {
		commandHandlers[k] = v
	}
	// Add Suggestions commands
	commands = append(commands, suggestions.Commands...)
	for k, v := range suggestions.CommandHandlers {
		commandHandlers[k] = v
	}
	// Add Report commands
	commands = append(commands, reports.Commands...)
	for k, v := range reports.CommandHandlers {
		commandHandlers[k] = v
	}

	var err error
	s, err = discordgo.New("Bot " + *shared.BotToken)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid bot parameters")
	}

	s.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsDirectMessages |
		discordgo.IntentsGuildMembers |
		discordgo.IntentMessageContent

	db.InitDB()

	if *shared.PrettyLogs {
		log.Logger = log.Output( //nolint:reassign // This only changes if the user prefers JSON over PrettyLogs
			zerolog.ConsoleWriter{Out: os.Stderr},
		)
	}

	// Add a handler for the interactionCreate event that will call the appropriate command handler
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			// Make sure it's an application command (e.g., /mycommand)
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
			return
		}
	})
	// Simply ready message when the bot is ready
	s.AddHandler(func(s *discordgo.Session, _ *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Open the session to begin listening for events
	err = s.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot open the session")
	}

	// Add interaction handlers, this might be messages, reactions etc... but not slash commands
	s.AddHandler(messageCreate)
	s.AddHandler(suggestions.UpvoteSuggestion)
	s.AddHandler(moderation.AlertHandler)

	// History handlers
	if *shared.MessageHistoryEnabled {
		s.AddHandler(history.OnMessageDelete)
		s.AddHandler(history.OnMessageUpdate)
	}
	if *shared.NicknameHistoryEnabled {
		s.AddHandler(history.OnGuildMemberUpdate)
	}

	// Register available slash commands
	log.Debug().Msg("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		var cmd *discordgo.ApplicationCommand
		cmd, err = s.ApplicationCommandCreate(s.State.User.ID, *shared.GuildID, v)
		if err != nil {
			log.Fatal().Err(err).Msgf("Cannot create '%v' command", v.Name)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	// This is used for e.g. the report command to know the guild name
	// This helps the user know where the report is coming from
	shared.SetGuildName(s)
	log.Debug().Msg("Guild name: " + shared.GuildName)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Info().Msg("Bot is now running. Press CTRL+C to exit.")
	<-stop
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.GuildID != "" { //nolint:nestif // This is required to know how to handle the message
		// Check if message is from a report channel
		// If it is, send the message to the user
		// If it isn't, ignore the message
		log.Debug().Msg("Incoming message from guild: " + m.GuildID)

		// Check if channel has category as parent
		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			log.Error().Err(err).Msg("Error fetching channel")
			return
		}

		if channel.ParentID == *shared.ReportCategory {
			userID, _ := db.GetReportByChannelID(m.ChannelID)
			log.Debug().Msg("Message is to user: " + userID)
			userChannel, usrChnlErr := s.UserChannelCreate(userID)
			if usrChnlErr != nil {
				log.Error().Err(usrChnlErr).Msg("Error creating user channel")
				return
			}

			// Send the users message to the channel
			_, err = s.ChannelMessageSend(userChannel.ID, m.Content)
			if err != nil {
				log.Error().Err(err).Msg("Error sending message")
				return
			}
		} else {
			db.CreateMessage(m.ID, m.Content, m.Author.ID, m.Author.Mention(), m.Timestamp.String(), m.ChannelID)
		}
		return
	}

	log.Debug().Msg("Incoming message from user:" + m.Author.ID)
	reports.OpenReport(s, m.ChannelID, m.Author.ID, m.Author.Username, m.Content, m.Author.AvatarURL("256x256"), false)
}
