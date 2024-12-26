package reports

import (
	"bytes"
	"context"
	"math280h/wisp/db"
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func Archive( //nolint:gocognit // This function is required to have multiple steps
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	// Archive the report
	// Delete the channel
	// Send a message to the user
	// Delete the command
	// Delete the command response
	currentChannel := i.ChannelID
	// If parent is the category, delete the channel
	channel, err := s.Channel(currentChannel)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch channel")
	}

	reportObj, err := shared.DBClient.Report.FindFirst(
		db.Report.ChannelID.Equals(currentChannel),
	).With(
		db.Report.User.Fetch(),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get report")
	}

	if channel.ParentID == *shared.ReportCategory { //nolint:nestif // This is required to know how to handle the message
		// Send all contents as a file to the archive channel
		messages, chlMsgErr := s.ChannelMessages(currentChannel, 100, "", "", "")
		if chlMsgErr != nil {
			log.Error().Err(chlMsgErr).Msg("Failed to fetch messages")
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
				author = reportObj.User().Nickname
			}
			file += author + ": " + message.Content + "\n"
		}

		// Send information about user
		reportedBy := "Reported by: " + messages[len(messages)-1].Author.Username + "\n"
		_, err = s.ChannelMessageSend(*shared.ArchiveChannel, reportedBy)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send message")
		}

		_, err = s.ChannelFileSend(*shared.ArchiveChannel, "report.txt", bytes.NewReader([]byte(file)))
		if err != nil {
			log.Error().Err(err).Msg("Failed to send file")
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
			log.Error().Err(err).Msg("Failed to send message")
		}

		// Send a message to the user
		userChannel, usrChnlErr := s.UserChannelCreate(i.Member.User.ID)
		if usrChnlErr != nil {
			log.Error().Err(usrChnlErr).Msg("Error creating user channel")
		}

		embed = &discordgo.MessageEmbed{
			Color:       0xff0000,
			Description: "Report archived",
		}
		_, err = s.ChannelMessageSendEmbed(userChannel.ID, embed)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send message")
		}

		// Delete the channel
		_, err = s.ChannelDelete(currentChannel)
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete channel")
		}

		// Close report in database
		_, err = shared.DBClient.Report.FindUnique(
			db.Report.ID.Equals(reportObj.ID),
		).Update(
			db.Report.Status.Set("closed"),
		).Exec(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("Failed to close report")
		}
	} else {
		// Send a ephemeral message to the user
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You can only use this command in a report channel",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to send fail interaction response for suggestion status")
		}
	}
}
