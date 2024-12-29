package moderation

import (
	"context"
	"errors"
	"math280h/wisp/db"
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func AddPointsToUser(userID string, pointsValue int) (int, bool) {
	userObj, err := shared.DBClient.User.FindFirst(
		db.User.UserID.Equals(userID),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user")
		return 0, false
	}

	userObj.Points += pointsValue

	// Update the user
	_, err = shared.DBClient.User.FindUnique(
		db.User.UserID.Equals(userID),
	).Update(
		db.User.Points.Set(userObj.Points),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user")
	}

	if userObj.Points >= *shared.MaxPoints {
		return userObj.Points, true
	}

	return userObj.Points, false
}

func AddInfraction(
	s *discordgo.Session,
	userID string,
	reason string,
	moderatorID string,
	moderatorName string,
	newPoints int,
	infractionType string,
) int {
	points, isOverLimit := AddPointsToUser(userID, newPoints)

	// Create the infraction
	_, err := shared.DBClient.Infraction.CreateOne(
		db.Infraction.User.Link(
			db.User.UserID.Equals(userID),
		),
		db.Infraction.Reason.Set(reason),
		db.Infraction.Type.Set(infractionType),
		db.Infraction.Points.Set(newPoints),
		db.Infraction.ModeratorID.Set(moderatorID),
		db.Infraction.ModeratorUsername.Set(moderatorName),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to create infraction")
	}

	if isOverLimit {
		err = s.GuildBanCreateWithReason(
			*shared.GuildID,
			userID,
			"User has exceeded the maximum amount of warning points",
			0,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to ban user")
		}
	}

	return points
}

func RemoveInfraction(s *discordgo.Session, i *discordgo.InteractionCreate) { //nolint:gocognit // Complexity is from response handling
	// Get the infraction ID
	infractionID := i.ApplicationCommandData().Options[0].IntValue()
	// Convert int64 to int
	infractionIDInt := int(infractionID)

	infraction, err := shared.DBClient.Infraction.FindUnique(
		db.Infraction.ID.Equals(infractionIDInt),
	).Delete().Exec(context.Background())
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Infraction not found",
					Flags:   64,
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to respond to remove infraction command")
			}
			return
		}

		log.Error().Err(err).Msg("Failed to get infraction to remove")
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something went wrong while trying to remove the infraction",
				Flags:   64,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to respond to remove infraction command")
		}
		return
	}

	// Remove the points from the user
	userObj, err := shared.DBClient.User.FindFirst(
		db.User.ID.Equals(infraction.UserID),
	).Exec(context.Background())
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "User not found",
					Flags:   64,
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to respond to remove infraction command")
			}
			return
		}
		log.Error().Err(err).Msg("Failed to get user")
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something went wrong while trying to remove the infraction",
				Flags:   64,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to respond to remove infraction command")
		}
	}

	// Check if the user is above the limit
	if userObj.Points >= *shared.MaxPoints {
		err = s.GuildBanDelete(*shared.GuildID, userObj.UserID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to unban user")
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to unban user",
					Flags:   64,
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to respond to remove infraction command")
			}
		}
	}

	userObj.Points -= infraction.Points

	// Update the user
	_, err = shared.DBClient.User.FindUnique(
		db.User.ID.Equals(infraction.UserID),
	).Update(
		db.User.Points.Set(userObj.Points),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user")
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something went wrong while trying to remove the infraction",
				Flags:   64,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to respond to remove infraction command")
		}
	}

	// Respond to the command
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Infraction removed",
			Flags:   64,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to respond to remove infraction command")
	}
}
