package suggestions

import (
	"math280h/wisp/internal/db"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func UpvoteSuggestion(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	// Get the suggestion ID
	var suggestionID = i.MessageComponentData().CustomID

	// ensure custom ID starts with vote_up or vote_down
	if !strings.HasPrefix(suggestionID, "vote_up") && !strings.HasPrefix(suggestionID, "vote_down") {
		return
	}

	// Split it by : to get the suggestion ID
	// The format is vote_up:<suggestion_id> or vote_down:<suggestion_id>
	var suggestionIDSplit = strings.Split(suggestionID, ":")
	if len(suggestionIDSplit) != 2 {
		return
	}

	// Get the vote type
	var voteType = suggestionIDSplit[0]
	// Get the suggestion ID
	suggestionIDStr, err := strconv.Atoi(suggestionIDSplit[1])
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert suggestion ID to int")
		return
	}
	existingSuggestionEmbed, _, _ := db.GetSuggestionByID(suggestionIDStr)
	if existingSuggestionEmbed == "" {
		log.Error().Msg("Suggestion not found")
		return
	}

	existingVoteID, existingVoteSentiment := db.GetSuggestionVoteByUserAndSuggestion(i.Member.User.ID, suggestionIDStr)
	switch voteType {
	case "vote_up":
		db.CreateSuggestionVote(suggestionIDStr, i.Member.User.ID, "upvote")
		if existingVoteID != 0 && existingVoteSentiment == "downvote" {
			db.DeleteSuggestionVote(existingVoteID)
		}
	case "vote_down":
		db.CreateSuggestionVote(suggestionIDStr, i.Member.User.ID, "downvote")
		if existingVoteID != 0 && existingVoteSentiment == "upvote" {
			db.DeleteSuggestionVote(existingVoteID)
		}
	}

	// Update the suggestion embed with the new upvote count
	// Get the upvote and downvote count
	upvotes, downvotes := db.GetSuggestionVoteCount(suggestionIDStr)
	embed := getSuggestionEmbed(suggestionIDStr, upvotes, downvotes, existingSuggestionEmbed, i.Member.User.ID)
	_, err = s.ChannelMessageEditEmbed(i.ChannelID, i.Message.ID, embed)
	if err != nil {
		log.Error().Err(err).Msg("Failed to edit suggestion embed")
	}

	// Send ephemeral message to the user that they have voted
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You have voted on the suggestion",
			Flags:   64,
		},
	})
	if err != nil {
		return
	}
}
