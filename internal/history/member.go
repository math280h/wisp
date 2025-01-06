package history

import (
	"context"
	"math280h/wisp/db"
	"math280h/wisp/internal/shared"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func updateNickIfChanged(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	userObj, err := shared.GetUserIfExists(&discordgo.User{
		ID: m.Member.User.ID,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user")
		return
	}

	newNick := m.Nick

	if m.Nick != "" && userObj.Nickname != newNick {
		log.Debug().Msg("Nickname changed: " + newNick)
		if newNick == "" {
			newNick = m.Member.User.Username
		}
		user := m.Member.User

		_, err = shared.DBClient.User.FindUnique(
			db.User.UserID.Equals(user.ID),
		).Update(
			db.User.Nickname.Set(newNick),
		).Exec(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("Failed to update nickname")
			return
		}

		embed := discordgo.MessageEmbed{
			Title:       "Nickname Changed",
			Description: user.Mention() + " (" + user.String() + ")",
			Color:       shared.DarkBlue,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Old Nickname",
					Value: userObj.Nickname,
				},
				{
					Name:  "New Nickname",
					Value: newNick,
				},
			},
		}

		_, channelErr := s.ChannelMessageSendEmbed(*shared.HistoryChannel, &embed)
		if channelErr != nil {
			log.Error().Err(channelErr).Msg("Failed to send message")
		}
	}
}

func OnGuildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	updateNickIfChanged(s, m)
}

func OnGuildMemeberJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	_, err := shared.DBClient.User.UpsertOne(
		db.User.UserID.Equals(m.Member.User.ID),
	).Create(
		db.User.UserID.Set(m.Member.User.ID),
		db.User.Nickname.Set(m.Member.User.Username),
		db.User.LastJoin.Set(time.Now().Format(time.RFC3339)),
	).Update(
		db.User.LastJoin.Set(time.Now().Format(time.RFC3339)),
		db.User.Nickname.Set(m.Member.User.Username),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return
	}

	embed := discordgo.MessageEmbed{
		Title:       "Member Joined",
		Description: m.Member.User.Mention() + " (" + m.Member.User.String() + ")",
		Color:       shared.Green,
	}

	_, channelErr := s.ChannelMessageSendEmbed(*shared.HistoryChannel, &embed)
	if channelErr != nil {
		log.Error().Err(channelErr).Msg("Failed to send message")
	}
}

func OnGuildMemberLeave(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	_, err := shared.DBClient.User.FindUnique(
		db.User.UserID.Equals(m.Member.User.ID),
	).Update(
		db.User.LastLeave.Set(time.Now().Format(time.RFC3339)),
		db.User.LeaveCount.Increment(1),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user")
		return
	}

	embed := discordgo.MessageEmbed{
		Title:       "Member Left",
		Description: m.Member.User.Mention() + " (" + m.Member.User.String() + ")",
		Color:       shared.Red,
	}

	_, channelErr := s.ChannelMessageSendEmbed(*shared.HistoryChannel, &embed)
	if channelErr != nil {
		log.Error().Err(channelErr).Msg("Failed to send message")
	}
}
