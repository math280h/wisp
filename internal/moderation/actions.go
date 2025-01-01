package moderation

import (
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func KickCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.ApplicationCommandData().Options[0].UserValue(s).ID
	reason := i.ApplicationCommandData().Options[1].StringValue()

	// Get user from the user ID
	user, err := s.User(userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get discord user")
		return
	}

	AddInfraction(
		s,
		userID,
		user.Username,
		reason,
		i.Member.User.ID,
		i.Member.User.Username,
		*shared.MaxPoints/2,
		"kick",
	)

	err = s.GuildMemberDeleteWithReason(*shared.GuildID, userID, reason)
	if err != nil {
		panic(err)
	}

	// Respond to the command
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User has been kicked",
		},
	})
	if err != nil {
		panic(err)
	}
}

func BanCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.ApplicationCommandData().Options[0].UserValue(s).ID
	reason := i.ApplicationCommandData().Options[1].StringValue()

	// Get user from the user ID
	user, err := s.User(userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get discord user")
		return
	}

	// Add 100% of the max points to the user
	AddInfraction(
		s,
		userID,
		user.Username,
		reason,
		i.Member.User.ID,
		i.Member.User.Username,
		*shared.MaxPoints,
		"ban",
	)

	// Respond to the command
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User has been banned",
		},
	})
	if err != nil {
		panic(err)
	}
}
