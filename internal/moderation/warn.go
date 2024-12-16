package moderation

import (
	"math280h/wisp/internal/db"
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
)

func addPointsToUser(userID string, pointsValue int) bool {
	rows, err := db.DBClient.Query("SELECT points FROM users WHERE user_id = ?", userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var points int
	for rows.Next() {
		err = rows.Scan(&points)
		if err != nil {
			panic(err)
		}
	}

	points += pointsValue
	_, err = db.DBClient.Exec("UPDATE users SET points = ? WHERE user_id = ?", points, userID)
	if err != nil {
		panic(err)
	}

	if points >= *shared.MaxPoints {
		return true
	}

	return false
}

func Warn(userID string, s *discordgo.Session) {
	isOverLimit := addPointsToUser(userID, *shared.WarnPoints)

	if isOverLimit {
		err := s.GuildBanCreateWithReason(*shared.GuildID, userID, "User has exceeded the maximum amount of warning points", 0)
		if err != nil {
			panic(err)
		}
	}
}

func Strike(userID string, s *discordgo.Session) {
	isOverLimit := addPointsToUser(userID, *shared.StrikePoints)

	if isOverLimit {
		err := s.GuildBanCreateWithReason(*shared.GuildID, userID, "User has exceeded the maximum amount of warning points", 0)
		if err != nil {
			panic(err)
		}
	}
}
