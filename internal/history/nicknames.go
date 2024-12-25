package history

import (
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func OnGuildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	var oldNick string
	var beforeUpdateFailed = false
	// TODO:: Improve this logic by keeping track of display name in user table
	if m.BeforeUpdate == nil {
		log.Debug().Msg("BeforeUpdate is nil, possibly because nick was changed while bot was booting up")
		oldNick = m.Member.User.String()
		beforeUpdateFailed = true
	} else {
		oldNick = m.BeforeUpdate.Nick
	}

	newNick := m.Nick

	if oldNick != newNick {
		user := m.Member.User

		var oldMsg string
		if beforeUpdateFailed {
			oldMsg = oldNick + " (Failed to get old nickname, possibly changed while bot was booting up)"
		} else {
			oldMsg = oldNick
		}

		embed := discordgo.MessageEmbed{
			Title:       "Nickname Changed",
			Description: user.Mention() + " (" + user.String() + ")",
			Color:       shared.DarkBlue,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Old Nickname",
					Value: oldMsg,
				},
				{
					Name:  "New Nickname",
					Value: newNick,
				},
			},
		}

		_, err := s.ChannelMessageSendEmbed(*shared.HistoryChannel, &embed)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send message")
		}
	}
}
