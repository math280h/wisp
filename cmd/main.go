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

func int64Ptr(i int64) *int64 {
	return &i
}

var (
	commands = []*discordgo.ApplicationCommand{ //nolint:gochecknoglobals // This is a list of commands for Discord
		{
			Name:                     "suggest",
			Description:              "Suggest something for the server",
			DefaultMemberPermissions: int64Ptr(discordgo.PermissionViewChannel),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "suggestion",
					Description: "The suggestion",
					Required:    true,
				},
			},
		},
		{
			Name:                     "suggestion-status",
			Description:              "Set the status of a suggestion",
			DefaultMemberPermissions: int64Ptr(discordgo.PermissionAdministrator),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "suggestion_id",
					Description: "The ID of the suggestion",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "status",
					Description: "The status of the suggestion",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Approved",
							Value: "approved",
						},
						{
							Name:  "Denied",
							Value: "denied",
						},
					},
				},
			},
		},
		{
			Name:                     "close",
			Description:              "Close the report",
			DefaultMemberPermissions: int64Ptr(discordgo.PermissionKickMembers),
		},
		{
			Name:                     "archive",
			Description:              "Archive the report",
			DefaultMemberPermissions: int64Ptr(discordgo.PermissionKickMembers),
		},
		{
			Name:                     "warn",
			Description:              "Warn the user",
			DefaultMemberPermissions: int64Ptr(discordgo.PermissionKickMembers),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to warn",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "The reason for the warning",
					Required:    true,
				},
			},
		},
		{
			Name:                     "strike",
			Description:              "Strike the user",
			DefaultMemberPermissions: int64Ptr(discordgo.PermissionKickMembers),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to strike",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "The reason for the strike",
					Required:    true,
				},
			},
		},
		{
			Name:                     "kick",
			Description:              "Kick the user",
			DefaultMemberPermissions: int64Ptr(discordgo.PermissionKickMembers),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to kick",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "The reason for the kick",
					Required:    true,
				},
			},
		},
		{
			Name:                     "ban",
			Description:              "Ban the user",
			DefaultMemberPermissions: int64Ptr(discordgo.PermissionBanMembers),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to ban",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "The reason for the ban",
					Required:    true,
				},
			},
		},
		{
			Name:                     "info",
			Description:              "Get information about a user",
			DefaultMemberPermissions: int64Ptr(discordgo.PermissionKickMembers),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to get information about",
					Required:    true,
				},
			},
		},
		{
			Name:                     "detain",
			Description:              "Detain the user",
			DefaultMemberPermissions: int64Ptr(discordgo.PermissionKickMembers),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to detain",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "The reason for the detainment",
					Required:    true,
				},
			},
		},
		{
			Name:                     "release",
			Description:              "Release the user",
			DefaultMemberPermissions: int64Ptr(discordgo.PermissionKickMembers),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to release",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "The reason for the release",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func( //nolint:gochecknoglobals // This is a map of commands to their handlers
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	){
		"close":             reports.Close,
		"archive":           reports.Archive,
		"warn":              moderation.WarnCommand,
		"strike":            moderation.StrikeCommand,
		"info":              moderation.InfoCommand,
		"suggest":           suggestions.CreateSuggestionCommand,
		"suggestion-status": suggestions.SetSuggestionStatusCommand,
		"detain":            moderation.DetainUserCommand,
		"release":           moderation.ReleaseUserCommand,
	}
)

func main() {
	shared.Init()

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

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			// Make sure it's an application command (e.g., /mycommand)
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
			return
		}
	})
	s.AddHandler(func(s *discordgo.Session, _ *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err = s.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot open the session")
	}

	s.AddHandler(messageCreate)
	s.AddHandler(suggestions.UpvoteSuggestion)
	s.AddHandler(moderation.AlertHandler)

	// History handlers
	s.AddHandler(history.OnGuildMemberUpdate)
	s.AddHandler(history.OnMessageDelete)
	s.AddHandler(history.OnMessageUpdate)

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
