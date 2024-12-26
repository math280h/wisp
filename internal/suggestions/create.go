package suggestions

import (
	"context"
	"math280h/wisp/db"
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

	suggestionObj, err := shared.DBClient.Suggestion.CreateOne(
		db.Suggestion.Suggestion.Set(suggestion),
		db.Suggestion.User.Link(
			db.User.UserID.Equals(i.Member.User.ID),
		),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to create suggestion")
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
	embed := getSuggestionEmbed(suggestionObj.ID, 0, 0, suggestion, i.Member.User.ID)

	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Up",
					Style:    discordgo.PrimaryButton,
					CustomID: "vote_up:" + strconv.Itoa(suggestionObj.ID),
					Emoji: &discordgo.ComponentEmoji{
						Name: "⬆️",
					},
				},
				discordgo.Button{
					Label:    "Down",
					Style:    discordgo.DangerButton,
					CustomID: "vote_down:" + strconv.Itoa(suggestionObj.ID),
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

	_, err = shared.DBClient.Suggestion.FindUnique(
		db.Suggestion.ID.Equals(suggestionObj.ID),
	).Update(
		db.Suggestion.EmbedID.Set(suggestionMessage.ID),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to set suggestion embed ID")
	}

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
