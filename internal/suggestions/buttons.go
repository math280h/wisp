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
	var suggestion_id string = i.MessageComponentData().CustomID

	// ensure custom ID starts with vote_up or vote_down
	if !strings.HasPrefix(suggestion_id, "vote_up") && !strings.HasPrefix(suggestion_id, "vote_down") {
		return
	}

	// Split it by : to get the suggestion ID
	// The format is vote_up:<suggestion_id> or vote_down:<suggestion_id>
	var suggestion_id_split []string = strings.Split(suggestion_id, ":")
	if len(suggestion_id_split) != 2 {
		return
	}

	// Get the vote type
	var vote_type string = suggestion_id_split[0]
	// Get the suggestion ID
	suggestion_id_str, err := strconv.Atoi(suggestion_id_split[1])
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert suggestion ID to int")
		return
	}
	existing_suggestion_embed, _, _ := db.GetSuggestionByID(suggestion_id_str)
	if existing_suggestion_embed == "" {
		log.Error().Msg("Suggestion not found")
		return
	}

	existing_vote_id, existing_vote_sentiment := db.GetSuggestionVoteByUserAndSuggestion(i.Member.User.ID, suggestion_id_str)
	switch vote_type {
	case "vote_up":
		db.CreateSuggestionVote(suggestion_id_str, i.Member.User.ID, "upvote")
		if existing_vote_id != 0 && existing_vote_sentiment == "downvote" {
			db.DeleteSuggestionVote(existing_vote_id)
		}
	case "vote_down":
		db.CreateSuggestionVote(suggestion_id_str, i.Member.User.ID, "downvote")
		if existing_vote_id != 0 && existing_vote_sentiment == "upvote" {
			db.DeleteSuggestionVote(existing_vote_id)
		}
	}

	// Update the suggestion embed with the new upvote count
	// Get the upvote and downvote count
	upvotes, downvotes := db.GetSuggestionVoteCount(suggestion_id_str)
	embed := getSuggestionEmbed(suggestion_id_str, upvotes, downvotes, existing_suggestion_embed, i.Member.User.ID)
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
