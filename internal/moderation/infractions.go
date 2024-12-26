package moderation

import (
	"context"
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
	points int,
	infractionType string,
) int {
	points, isOverLimit := AddPointsToUser(userID, points)

	// Create the infraction
	_, err := shared.DBClient.Infraction.CreateOne(
		db.Infraction.User.Link(
			db.User.UserID.Equals(userID),
		),
		db.Infraction.Reason.Set(reason),
		db.Infraction.Type.Set(infractionType),
		db.Infraction.Points.Set(points),
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
