package core

import (
	"math280h/wisp/internal/db"
	"math280h/wisp/internal/reports"
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func HandleIncomingMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
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
