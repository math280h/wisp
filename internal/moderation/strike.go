package moderation

import (
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
)

func StrikeCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.ApplicationCommandData().Options[0].UserValue(s).ID
	reason := i.ApplicationCommandData().Options[1].StringValue()
	points := AddInfraction(s, userID, reason, i.Member.User.ID, *shared.StrikePoints)

	// Respond to the command
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User has been striked",
		},
	})
	if err != nil {
		panic(err)
	}

	// Inform the user that they have been striked
	embed := GenerateInfractionEmbed(reason, points)

	userChannel, err := s.UserChannelCreate(userID)
	if err != nil {
		panic(err)
	}

	_, err = s.ChannelMessageSendEmbed(userChannel.ID, embed)
	if err != nil {
		panic(err)
	}
}
