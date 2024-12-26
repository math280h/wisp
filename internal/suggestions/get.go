package suggestions

import (
	"math280h/wisp/internal/shared"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func getSuggestionEmbed(id int, upvotes int, downvotes int, suggestion string, author string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       shared.DarkBlue,
		Title:       "New Suggestion!",
		Description: "A new suggestion has been created, you can read it below and vote on it.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Created By",
				Value: "<@" + author + ">",
			},
			{
				Name:   "Suggestion",
				Value:  suggestion,
				Inline: true,
			},
			// Upvote and Downvote Count
			{
				Name:   "Upvotes",
				Value:  strconv.Itoa(upvotes),
				Inline: true,
			},
			{
				Name:   "Downvotes",
				Value:  strconv.Itoa(downvotes),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Suggestion ID: " + strconv.Itoa(id),
		},
	}
}
