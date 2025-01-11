package moderation

import (
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func WarnCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.ApplicationCommandData().Options[0].UserValue(s).ID
	reason := i.ApplicationCommandData().Options[1].StringValue()

	// Get user from the user ID
	user, err := s.User(userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get discord user")
		return
	}

	points := AddInfraction(
		s,
		userID,
		user.Username,
		reason,
		i.Member.User.ID,
		i.Member.User.Username,
		*shared.WarnPoints,
		"warn",
	)

	// Respond to the command
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User has been warned",
			Flags:   64,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to respond to warn command")
	}

	// Inform the user that they have been warned
	embed := GenerateInfractionEmbed(reason, points)

	userChannel, err := s.UserChannelCreate(userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user channel")
	}

	_, err = s.ChannelMessageSendEmbed(userChannel.ID, embed)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send warn message")
	}
}
