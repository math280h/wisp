package history

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"

	"math280h/wisp/internal/db"
	"math280h/wisp/internal/shared"
)

type MessageData struct {
	ID        string
	Content   string
	AuthorID  string
	AuthorTag string
	Timestamp string
}

func OnMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	// Check if the message content was changed
	content, _, authorTag, _, _, _ := db.GetMessageByID(m.ID)

	if content != "" && content != m.Content {
		embed := discordgo.MessageEmbed{
			Title:       "Message Edited",
			Description: "A message by " + authorTag + " was edited.",
			Color:       shared.Orange,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Old Content",
					Value: content,
				},
				{
					Name:  "New Content",
					Value: m.Content,
				},
				{
					Name: "Link",
					Value: "[Jump to message](https://discord.com/channels/" +
						m.GuildID + "/" + m.ChannelID + "/" + m.ID + ")",
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Message ID: " + m.ID,
			},
		}

		db.UpdateContentByID(m.ID, m.Content)

		_, err := s.ChannelMessageSendEmbed(*shared.HistoryChannel, &embed)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send message edit embed")
		}
	}
}

func OnMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	content, _, authorTag, _, _, _ := db.GetMessageByID(m.ID)

	embed := discordgo.MessageEmbed{
		Title:       "Message Deleted",
		Description: "A message by " + authorTag + " was deleted.",
		Color:       shared.Red,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Content",
				Value: content,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Message ID: " + m.ID,
		},
	}

	db.DeleteMessageByID(m.ID)

	_, err := s.ChannelMessageSendEmbed(*shared.HistoryChannel, &embed)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message delete embed")
	}
}
