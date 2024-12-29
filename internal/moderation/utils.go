package moderation

import (
	"context"
	"math280h/wisp/db"
	"math280h/wisp/internal/shared"
	"strconv"

	"github.com/rs/zerolog/log"
)

func getUserReportCount(userID int) int {
	var res []struct {
		Reports db.String `json:"reports"`
	}
	err := shared.DBClient.Prisma.QueryRaw("SELECT COUNT(id) as reports FROM reports WHERE user_id = ?", userID).
		Exec(
			context.Background(),
			&res,
		)

	if err != nil {
		log.Error().Err(err).Msg("Failed to send fail interaction response for info")
	}

	// Convert the reports to an int
	reports, err := strconv.Atoi(res[0].Reports)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert reports to int")
	}
	return reports
}
