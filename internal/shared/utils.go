package shared

import (
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func Int64Ptr(i int64) *int64 {
	return &i
}

func StringTimeToDiscordTimestamp(t string) string {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", t)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse time")
	}
	unixTimestamp := parsedTime.Unix()
	return "<t:" + strconv.FormatInt(unixTimestamp, 10) + ":R>"
}

func StringWithTzToDiscordTimestamp(t string) string {
	if t == "Never" {
		return "Never"
	}

	parsedTime, err := time.Parse("2006-01-02T15:04:05-07:00", t)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse time")
	}
	unixTimestamp := parsedTime.Unix()
	return "<t:" + strconv.FormatInt(unixTimestamp, 10) + ":R>"
}
