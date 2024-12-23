package suggestions

import (
	"math280h/wisp/internal/db"
	"math280h/wisp/internal/shared"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
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

func CreateSuggestionCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Ensure suggestion came from shared.SuggestionsChannel
	if i.ChannelID != *shared.SuggestionChannel {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command can only be used in the suggestions channel",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to send fail interaction response for suggestion")
		}
		return
	}

	suggestion := i.ApplicationCommandData().Options[0].StringValue()

	// Inform the user that their suggestion has been created
	suggestionID := db.CreateSuggestion(i.Member.User.ID, suggestion)
	embed := getSuggestionEmbed(suggestionID, 0, 0, suggestion, i.Member.User.ID)

	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Up",
					Style:    discordgo.PrimaryButton,
					CustomID: "vote_up:" + strconv.Itoa(suggestionID),
					Emoji: &discordgo.ComponentEmoji{
						Name: "⬆️",
					},
				},
				discordgo.Button{
					Label:    "Down",
					Style:    discordgo.DangerButton,
					CustomID: "vote_down:" + strconv.Itoa(suggestionID),
					Emoji: &discordgo.ComponentEmoji{
						Name: "⬇️",
					},
				},
			},
		},
	}

	// Send to shared.SuggestionsChannel
	suggestionMessage, err := s.ChannelMessageSendComplex(*shared.SuggestionChannel, &discordgo.MessageSend{
		Embed:      embed,
		Components: buttons,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send suggestion message")
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create suggestion",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to send fail interaction response for suggestion")
		}
		return
	}

	db.SetSuggestionEmbedID(suggestionID, suggestionMessage.ID)

	// Respond to the command
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Suggestion has been created",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		panic(err)
	}
}

func SetSuggestionStatusCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get suggestion ID
	suggestionID := i.ApplicationCommandData().Options[0].IntValue()
	// Convert int64 to int
	suggestionIDInt := int(suggestionID)
	// Get status
	status := i.ApplicationCommandData().Options[1].StringValue()

	// Ensure status is either approved or denied
	if status != "approved" && status != "denied" {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Status must be either 'approved' or 'denied'",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to send fail interaction response for suggestion status")
		}
		return
	}

	// Get the suggestion embed
	embedID, suggestionUser, suggestionMessage := db.GetSuggestionByID(suggestionIDInt)

	db.SetSuggestionStatus(suggestionIDInt, status)

	// Set color and title
	var color int
	var title string
	switch status {
	case "approved":
		color = shared.Green
		title = "Approved Suggestion!"
	case "denied":
		color = shared.Red
		title = "Denied Suggestion!"
	}

	// Update the suggestion embed with the new status, and remove the buttons
	embed := &discordgo.MessageEmbed{
		Color:       color,
		Title:       title,
		Description: "The suggestion has been " + status,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Created By",
				Value: "<@" + suggestionUser + ">",
			},
			{
				Name:   "Suggestion",
				Value:  suggestionMessage,
				Inline: true,
			},
		},
	}

	// Remove buttons from the suggestion
	_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    *shared.SuggestionChannel,
		ID:         embedID,
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Components: &[]discordgo.MessageComponent{},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to edit suggestion embed")
	}

	// Respond to the command
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Suggestion status has been set",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send success interaction response for suggestion status")
	}
}
