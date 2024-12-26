package moderation

import (
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
)

func KickCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.ApplicationCommandData().Options[0].UserValue(s).ID
	reason := i.ApplicationCommandData().Options[1].StringValue()

	AddInfraction(
		s,
		userID,
		reason,
		i.Member.User.ID,
		i.Member.User.Username,
		*shared.MaxPoints/2,
		"kick",
	)

	err := s.GuildMemberDeleteWithReason(*shared.GuildID, userID, reason)
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

	// Add 100% of the max points to the user
	AddInfraction(
		s,
		userID,
		reason,
		i.Member.User.ID,
		i.Member.User.Username,
		*shared.MaxPoints,
		"ban",
	)

	// Respond to the command
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User has been banned",
		},
	})
	if err != nil {
		panic(err)
	}
}
