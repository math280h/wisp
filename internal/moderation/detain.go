package moderation

import (
	"math280h/wisp/internal/reports"
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func DetainUserCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	target := i.ApplicationCommandData().Options[0].UserValue(s)
	reason := i.ApplicationCommandData().Options[1].StringValue()
	// Give user the muted role
	err := s.GuildMemberRoleAdd(*shared.GuildID, target.ID, *shared.MutedRole)
	if err != nil {
		log.Error().Err(err).Msg("Error adding role to user")
	}

	// Respond to the command
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User has been detained",
			Flags:   64,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Error responding to command")
	}

	// Inform the user that they have been detained

	userChannel, err := s.UserChannelCreate(target.ID)
	if err != nil {
		log.Error().Err(err).Msg("Error creating user channel")
	}

	_, err = s.ChannelMessageSend(userChannel.ID, "You have been detained for the following reason: "+reason)
	if err != nil {
		log.Error().Err(err).Msg("Error sending message")
	}

	// reports.OpenReport(s, m.ChannelID, m.Author.ID, m.Author.Username, m.Content, m.Author.AvatarURL("256x256"))
	reports.OpenReport(
		s,
		userChannel.ID,
		target.ID,
		target.Username,
		reason,
		target.AvatarURL("256x256"),
		true,
	)
}

func ReleaseUserCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	target := i.ApplicationCommandData().Options[0].UserValue(s)
	reason := i.ApplicationCommandData().Options[1].StringValue()
	// Remove the muted role from the user
	err := s.GuildMemberRoleRemove(*shared.GuildID, target.ID, *shared.MutedRole)
	if err != nil {
		log.Error().Err(err).Msg("Error removing role from user")
	}

	// Respond to the command
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User has been released",
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Error responding to command")
	}

	// Inform the user that they have been released

	userChannel, err := s.UserChannelCreate(target.ID)
	if err != nil {
		log.Error().Err(err).Msg("Error creating user channel")
	}

	_, err = s.ChannelMessageSend(userChannel.ID, "You have been released for the following reason: "+reason)
	if err != nil {
		log.Error().Err(err).Msg("Error sending message")
	}
}
