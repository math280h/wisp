package history

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"

	"math280h/wisp/db"
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
	// Ensure the message is not from the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if the message content was changed
	messageObj, err := shared.DBClient.Message.FindFirst(
		db.Message.ID.Equals(m.ID),
	).With(
		db.Message.Author.Fetch(),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get message")
		return
	}

	if messageObj.Content != "" && messageObj.Content != m.Content {
		embed := discordgo.MessageEmbed{
			Title:       "Message Edited",
			Description: "A message by <@" + messageObj.Author().UserID + "> was edited.",
			Color:       shared.Orange,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Old Content",
					Value: messageObj.Content,
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

		_, err = shared.DBClient.Message.FindUnique(
			db.Message.ID.Equals(messageObj.ID),
		).Update(
			db.Message.Content.Set(m.Content),
		).Exec(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("Failed to update message content")
		}

		_, chnlErr := s.ChannelMessageSendEmbed(*shared.HistoryChannel, &embed)
		if chnlErr != nil {
			log.Error().Err(err).Msg("Failed to send message edit embed")
		}
	}
}

func OnMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	messageObj, err := shared.DBClient.Message.FindFirst(
		db.Message.ID.Equals(m.ID),
	).With(
		db.Message.Author.Fetch(),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get message")
		return
	}

	embed := discordgo.MessageEmbed{
		Title:       "Message Deleted",
		Description: "A message by <@" + messageObj.Author().UserID + "> was deleted.",
		Color:       shared.Red,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Content",
				Value: messageObj.Content,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Message ID: " + m.ID,
		},
	}

	_, err = shared.DBClient.Message.FindUnique(
		db.Message.ID.Equals(messageObj.ID),
	).Delete().Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete message")
	}

	_, err = s.ChannelMessageSendEmbed(*shared.HistoryChannel, &embed)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message delete embed")
	}
}
