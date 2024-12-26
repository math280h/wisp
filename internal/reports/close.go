package reports

import (
	"math280h/wisp/internal/db"
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("Failed to fetch channel")
	}

	if channel.ParentID == *shared.ReportCategory { //nolint:nestif // This is required to know how to handle the message
		// Get the report
		userID, _ := db.GetReportByChannelID(currentChannel)

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
		_, err = s.ChannelMessageSendEmbed(*shared.LogChannel, embed)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send message")
		}

		// Send a message to the user
		userChannel, usrChnlErr := s.UserChannelCreate(userID)
		if usrChnlErr != nil {
			log.Error().Err(usrChnlErr).Msg("Error creating user channel")
		}

		embed = &discordgo.MessageEmbed{
			Color:       0xff0000,
			Description: "Report closed",
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
		db.CloseReport(currentChannel)
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
