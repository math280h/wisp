package moderation

import (
	"context"
	"errors"
	"math280h/wisp/db"
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func NoteCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.ApplicationCommandData().Options[0].UserValue(s).ID
	text := i.ApplicationCommandData().Options[1].StringValue()

	// Get user from the user ID
	user, err := s.User(userID)
	if err != nil {
		shared.SimpleEphemeralInteractionResponse("Failed to get discord user", s, i.Interaction)
		log.Error().Err(err).Msg("Failed to get discord user")
		return
	}

	userObj, userErr := shared.GetUserIfExists(user)
	if userErr != nil {
		if errors.Is(userErr, db.ErrNotFound) {
			shared.SimpleEphemeralInteractionResponse("User not found", s, i.Interaction)
			return
		}
		shared.SimpleEphemeralInteractionResponse("Failed to get or create user", s, i.Interaction)
		log.Error().Err(userErr).Msg("Failed to get or create user")
		return
	}

	_, noteErr := shared.DBClient.Note.CreateOne(
		db.Note.User.Link(
			db.User.UserID.Equals(userObj.UserID),
		),
		db.Note.Content.Set(text),
		db.Note.ModeratorID.Set(i.Member.User.ID),
		db.Note.ModeratorUsername.Set(i.Member.User.Username),
	).Exec(context.Background())
	if noteErr != nil {
		shared.SimpleEphemeralInteractionResponse("Failed to create note", s, i.Interaction)
		log.Error().Err(noteErr).Msg("Failed to create note")
		return
	}

	// Respond to the command
	shared.SimpleEphemeralInteractionResponse("Note added to user", s, i.Interaction)
}
