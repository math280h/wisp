package moderation

import (
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func AlertHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// If the react is :stop_sign:
	if r.Emoji.Name == "🛑" {
		log.Debug().Msg("User reacted with 🛑")

		embed := &discordgo.MessageEmbed{
			Color:       0xe74c3c,
			Title:       "Alert",
			Description: "A user has reacted with 🛑 to this message",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Channel",
					Value:  r.ChannelID,
					Inline: true,
				},
				{
					Name:   "Message ID",
					Value:  r.MessageID,
					Inline: true,
				},
				{
					Name:   "Reported By",
					Value:  "<@" + r.UserID + ">",
					Inline: true,
				},
				{
					Name:  "Message Link",
					Value: "https://discord.com/channels/" + *shared.GuildID + "/" + r.ChannelID + "/" + r.MessageID,
				},
			},
		}

		_, err := s.ChannelMessageSendEmbed(*shared.AlertChannel, embed)
		if err != nil {
			log.Error().Err(err).Msg("Error sending message")
			return
		}
	}
}
