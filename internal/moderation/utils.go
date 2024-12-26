package moderation

import (
	"context"
	"math280h/wisp/db"
	"math280h/wisp/internal/shared"

	"github.com/rs/zerolog/log"
)

func getUserReportCount(userID int) int {
	var res []struct {
		Reports db.RawInt `json:"reports"`
	}
	err := shared.DBClient.Prisma.QueryRaw("SELECT COUNT(*) as reports FROM reports WHERE user_id = ?", userID).
		Exec(
			context.Background(),
			&res,
		)

	if err != nil {
		log.Error().Err(err).Msg("Failed to send fail interaction response for info")
	}
	// Convert the reports to an int
	return int(res[0].Reports)
}
