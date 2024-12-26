package moderation

import (
	"math280h/wisp/internal/shared"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func GenerateInfractionEmbed(reason string, points int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       0xe74c3c,
		Title:       "You have been warned",
		Description: "You have been warned by a moderator. Please ensure you follow the rules in the future.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Guild",
				Value: shared.GuildName,
			},
			{
				Name:   "Reason",
				Value:  reason,
				Inline: true,
			},
			{
				Name:   "Total Points",
				Value:  strconv.Itoa(points) + " / " + strconv.Itoa(*shared.MaxPoints),
				Inline: true,
			},
		},
	}
}
