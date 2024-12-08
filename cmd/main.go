package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"

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
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	db.InitDb()
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
		log.Fatalf("Cannot open the session: %v", err)
	}

	s.AddHandler(messageCreate)
	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *shared.GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	s.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
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
		fmt.Println("Message is from guild:", m.GuildID)

		// Check if channel has category as parent
		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			fmt.Println("Error fetching channel:", err)
		}

		if channel.ParentID == *shared.ReportCategory {
			fmt.Println("Message is from report channel")
			// TODO:: Store this in db
			userId := "128897567349145600"
			fmt.Println("User ID Created:", userId)
			userChannel, err := s.UserChannelCreate(userId)
			if err != nil {
				fmt.Println("Error creating user channel:", err)
				return
			}

			// Send the users message to the channel
			_, err = s.ChannelMessageSend(userChannel.ID, m.Content)
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		} else {
			fmt.Println("Message is not from report channel")
		}

		return
	}

	fmt.Println("Got message from", m.Author.ID)
	if m.Author.ID == s.State.User.ID {
		fmt.Println("message is from bot")
		return
	}

	// Check if the channel exists in the specified category
	expectedChannelName := m.Author.Username
	// Remove any special characters from the channel name
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	expectedChannelName = re.ReplaceAllString(expectedChannelName, "")

	channelExists := false
	channelId, channelName := db.GetReportByUserID(m.Author.ID)
	if channelName != "" {
		channelExists = true
	}

	if !channelExists {
		// Create a new channel in the specified category
		channelData := &discordgo.GuildChannelCreateData{
			Name:     expectedChannelName,
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: *shared.ReportCategory,
		}

		newChannel, err := s.GuildChannelCreateComplex(*shared.GuildID, *channelData)
		if err != nil {
			fmt.Println("Error creating channel:", err)
			return
		}
		fmt.Println("Created channel:", newChannel)
		db.CreateReport(newChannel.ID, newChannel.Name, m.Author.ID)

		// Send embed message to the channel
		embed := &discordgo.MessageEmbed{
			Color: 0x00ff00,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "User (Tag)",
					Value:  "<@" + m.Author.ID + ">",
					Inline: true,
				},
				{
					Name:   "User (Username)",
					Value:  m.Author.Username,
					Inline: true,
				},
				{
					Name:   "User (ID)",
					Value:  m.Author.ID,
					Inline: false,
				},
			},
			// Set the image as the users avatar
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: m.Author.AvatarURL("256x256"),
			},
		}

		// Send an embed to the user saying that the channel has been created
		user_embed := &discordgo.MessageEmbed{
			Color:       0x00ff00,
			Title:       "New report opened",
			Description: "A new report has been opened for you. Please use this channel to communicate with the staff.",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Channel",
					Value:  "<#" + newChannel.ID + ">",
					Inline: true,
				},
			},
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, user_embed)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		_, err = s.ChannelMessageSendEmbed(newChannel.ID, embed)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		// Send the users message to the channel
		_, err = s.ChannelMessageSend(newChannel.ID, m.Content)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	} else {
		fmt.Println("Open channel found, sending message to channel")
		// Send the users message to the channel
		_, err := s.ChannelMessageSend(channelId, m.Content)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}

	// Create a new channel in the
}
