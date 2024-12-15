package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"math280h/wisp/internal/db"
	"math280h/wisp/internal/reports"
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

var s *discordgo.Session

func init() { shared.Init() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *shared.BotToken)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid bot parameters")
	}

	db.InitDb()

	if *shared.PrettyLogs {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "close",
			Description: "Close the report",
		},
		{
			Name:        "archive",
			Description: "Archive the report",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"close":   reports.Close,
		"archive": reports.Archive,
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot open the session")
	}

	s.AddHandler(messageCreate)
	log.Debug().Msg("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *shared.GuildID, v)
		if err != nil {
			log.Fatal().Err(err).Msgf("Cannot create '%v' command", v.Name)
		}
		registeredCommands[i] = cmd
	}

	s.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	defer s.Close()

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

	if m.GuildID != "" {
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

			userId, _ := db.GetReportByChannelID(m.ChannelID)
			log.Debug().Msg("Message is to user: " + userId)
			userChannel, err := s.UserChannelCreate(userId)
			if err != nil {
				log.Error().Err(err).Msg("Error creating user channel")
				return
			}

			// Send the users message to the channel
			_, err = s.ChannelMessageSend(userChannel.ID, m.Content)
			if err != nil {
				log.Error().Err(err).Msg("Error sending message")
				return
			}
		}
		return
	}

	log.Debug().Msg("Incoming message from user:" + m.Author.ID)

	// Check if the channel exists in the specified category
	expectedChannelName := m.Author.Username
	// Remove any special characters from the channel name
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	expectedChannelName = re.ReplaceAllString(expectedChannelName, "")

	db.CreateUserIfNotExist(m.Author.ID, m.Author.Username)

	channelExists := false
	channelId, channelName := db.GetReportByUserID(m.Author.ID)
	if channelName != "" {
		channelExists = true
	}

	if !channelExists {
		// Create a new channel in the specified category
		_, err := reports.CreateReportChannel(expectedChannelName, s, m)
		if err != nil {
			_, err = s.ChannelMessageSend(m.ChannelID, "There was an error sending your message. Please try again.")
			if err != nil {
				log.Error().Err(err).Msg("Error sending message to channel " + m.ChannelID)
			}
		}
	} else {
		log.Debug().Msg("Open channel found, sending message to channel")
		// Send the users message to the channel
		_, err := s.ChannelMessageSend(channelId, m.Content)
		if err != nil {
			log.Error().Err(err).Msg("Error sending message to channel " + channelId)

			// If the error message contains Unknown Channel, attempt to create a new channel
			if fmt.Sprint(err) == "HTTP 404 Not Found, {\"message\": \"Unknown Channel\", \"code\": 10003}" {
				log.Debug().Msg("Channel not found, attempting to creating new channel")

				_, err := reports.CreateReportChannel(expectedChannelName, s, m)
				if err != nil {
					_, err = s.ChannelMessageSend(m.ChannelID, "There was an error sending your message. Please try again.")
					if err != nil {
						log.Error().Err(err).Msg("Error sending message to channel " + m.ChannelID)
					}
				}
			}
			return
		}
	}
}
