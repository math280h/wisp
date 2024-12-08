package reports

import (
	"bytes"
	"fmt"
	"math280h/wisp/internal/db"
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
)

func Close(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Close the report
	// Delete the channel
	// Send a message to the user
	// Delete the command
	// Delete the command response
	currentChannel := i.ChannelID
	// If parent is the category, delete the channel
	channel, err := s.Channel(currentChannel)
	if err != nil {
		fmt.Println("Error fetching channel:", err)
	}

	if channel.ParentID == *shared.ReportCategory {
		// Send a message to report log channel
		embed := &discordgo.MessageEmbed{
			Color:       0xff0000,
			Description: "Report closed",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Report",
					Value:  channel.Name,
					Inline: true,
				},
				{
					Name:   "Closed by",
					Value:  "<@" + i.Member.User.ID + ">",
					Inline: true,
				},
			},
		}
		_, err := s.ChannelMessageSendEmbed(*shared.LogChannel, embed)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}

		// Send a message to the user
		userChannel, err := s.UserChannelCreate(i.Member.User.ID)
		if err != nil {
			fmt.Println("Error creating user channel:", err)
		}

		embed = &discordgo.MessageEmbed{
			Color:       0xff0000,
			Description: "Report closed",
		}
		_, err = s.ChannelMessageSendEmbed(userChannel.ID, embed)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}

		// Delete the channel
		_, err = s.ChannelDelete(currentChannel)
		if err != nil {
			fmt.Println("Error deleting channel:", err)
		}

		// Close report in database
		db.CloseReport(currentChannel)
	} else {
		// Send a ephemeral message to the user
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You can only use this command in a report channel",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}

func Archive(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Archive the report
	// Delete the channel
	// Send a message to the user
	// Delete the command
	// Delete the command response
	currentChannel := i.ChannelID
	// If parent is the category, delete the channel
	channel, err := s.Channel(currentChannel)
	if err != nil {
		fmt.Println("Error fetching channel:", err)
	}
	_, reportChannelName := db.GetReportByChannelID(currentChannel)

	if channel.ParentID == *shared.ReportCategory {
		// Send all contents as a file to the archive channel
		messages, err := s.ChannelMessages(currentChannel, 100, "", "", "")
		if err != nil {
			fmt.Println("Error fetching messages:", err)
		}

		// Create a file with all the messages
		file := ""
		// Iterate in reverse order
		for i := len(messages) - 1; i >= 0; i-- {
			// Skip first message
			if i == len(messages)-1 {
				continue
			}
			message := messages[i]
			// If author is a bot, replace the name
			author := message.Author.Username
			if message.Author.Bot {
				author = reportChannelName
			}
			file += author + ": " + message.Content + "\n"
		}

		// Send information about user
		reportedBy := "Reported by: " + messages[len(messages)-1].Author.Username + "\n"
		_, err = s.ChannelMessageSend(*shared.ArchiveChannel, reportedBy)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}

		_, err = s.ChannelFileSend(*shared.ArchiveChannel, "report.txt", bytes.NewReader([]byte(file)))
		if err != nil {
			fmt.Println("Error sending file:", err)
		}

		logChannel := "1281803878278500352"
		// Send a message to report log channel
		embed := &discordgo.MessageEmbed{
			Color:       0xff0000,
			Description: "Report archived",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Report",
					Value:  channel.Name,
					Inline: true,
				},
				{
					Name:   "Archived by",
					Value:  "<@" + i.Member.User.ID + ">",
					Inline: true,
				},
			},
		}
		_, err = s.ChannelMessageSendEmbed(logChannel, embed)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}

		// Send a message to the user
		userChannel, err := s.UserChannelCreate(i.Member.User.ID)
		if err != nil {
			fmt.Println("Error creating user channel:", err)
		}

		embed = &discordgo.MessageEmbed{
			Color:       0xff0000,
			Description: "Report archived",
		}
		_, err = s.ChannelMessageSendEmbed(userChannel.ID, embed)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}

		// Delete the channel
		_, err = s.ChannelDelete(currentChannel)
		if err != nil {
			fmt.Println("Error deleting channel:", err)
		}

		// Close report in database
		db.CloseReport(currentChannel)
	} else {
		// Send a ephemeral message to the user
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You can only use this command in a report channel",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}
