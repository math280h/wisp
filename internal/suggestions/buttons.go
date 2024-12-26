package suggestions

import (
	"context"
	"math280h/wisp/db"
	"math280h/wisp/internal/shared"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func SuggestionVote(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	existingSuggestionEmbed, err := shared.DBClient.Suggestion.FindFirst(
		db.Suggestion.ID.Equals(suggestionIDStr),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get suggestion")
		return
	}

	// Get vote by user and suggestion
	userObj, err := shared.DBClient.User.UpsertOne(
		db.User.UserID.Equals(i.Member.User.ID),
	).Create(
		db.User.UserID.Set(i.Member.User.ID),
		db.User.Nickname.Set(i.Member.User.Username),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user")
		return
	}

	existingVote, err := shared.DBClient.SuggestionVote.FindFirst(
		db.SuggestionVote.SuggestionID.Equals(suggestionIDStr),
		db.SuggestionVote.UserID.Equals(userObj.ID),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get existing vote")
		return
	}

	if existingVote != nil && existingVote.Sentiment != voteType {
		_, err = shared.DBClient.SuggestionVote.FindUnique(
			db.SuggestionVote.ID.Equals(existingVote.ID),
		).Delete().Exec(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete existing vote")
			return
		}

		_, err = shared.DBClient.SuggestionVote.CreateOne(
			db.SuggestionVote.Suggestion.Link(
				db.Suggestion.ID.Equals(suggestionIDStr),
			),
			db.SuggestionVote.User.Link(
				db.User.ID.Equals(userObj.ID),
			),
			db.SuggestionVote.Sentiment.Set(voteType),
		).Exec(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("Failed to create new vote")
			return
		}
	} else {
		log.Debug().Msg("Vote ignored as duplicate")
		return
	}

	// Update the suggestion embed with the new upvote count
	// Get the upvote and downvote count
	var res []struct {
		Upvotes   db.RawInt `json:"upvotes"`
		Downvotes db.RawInt `json:"downvotes"`
	}
	err = shared.DBClient.Prisma.QueryRaw(
		`SELECT 
        	SUM(CASE WHEN sentiment = 'vote_up' THEN 1 ELSE 0 END) AS upvotes,
        	SUM(CASE WHEN sentiment = 'vote_down' THEN 1 ELSE 0 END) AS downvotes
     	FROM suggestion_votes
     	WHERE suggestion_id = ?`,
		existingVote.SuggestionID,
	).Exec(context.Background(), &res)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get suggestion vote count")
		return
	}
	// Convert the reports to an int
	upvotes := int(res[0].Upvotes)
	downvotes := int(res[0].Downvotes)

	// Ensure EmbedID is not nil
	embedID, ok := existingSuggestionEmbed.EmbedID()
	if !ok {
		log.Error().Msg("Failed to get embed ID")
	}

	embed := getSuggestionEmbed(existingVote.SuggestionID, upvotes, downvotes, embedID, i.Member.User.ID)
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
