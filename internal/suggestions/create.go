package suggestions

import (
	"math280h/wisp/internal/db"
	"math280h/wisp/internal/shared"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

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
