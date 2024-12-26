package suggestions

import (
	"math280h/wisp/internal/db"
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

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
