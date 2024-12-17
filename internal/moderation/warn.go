package moderation

import (
	"math280h/wisp/internal/db"
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
)

func addPointsToUser(userID string, pointsValue int) (current_points int, over bool) {
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
		return points, true
	}

	return points, false
}

func addWarnEntry(userID string, reason string, moderatorID string) {
	_, err := db.DBClient.Exec("INSERT INTO warns (user_id, reason, moderator_id) VALUES (?, ?, ?)", userID, reason, moderatorID)
	if err != nil {
		panic(err)
	}
}

func Warn(userID string, reason string, moderator_id string, s *discordgo.Session) int {
	points, isOverLimit := addPointsToUser(userID, *shared.WarnPoints)
	addWarnEntry(userID, reason, moderator_id)

	if isOverLimit {
		err := s.GuildBanCreateWithReason(*shared.GuildID, userID, "User has exceeded the maximum amount of warning points", 0)
		if err != nil {
			panic(err)
		}
	}

	return points
}

func Strike(userID string, reason string, moderator_id string, s *discordgo.Session) int {
	points, isOverLimit := addPointsToUser(userID, *shared.StrikePoints)
	addWarnEntry(userID, reason, moderator_id)

	if isOverLimit {
		err := s.GuildBanCreateWithReason(*shared.GuildID, userID, "User has exceeded the maximum amount of warning points", 0)
		if err != nil {
			panic(err)
		}
	}

	return points
}
