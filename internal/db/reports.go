package db

import (
	"math/rand"
	"strconv"
)

func GetReportByUserID(userID string) (string, string) {
	rows, err := DBClient.Query(
		"SELECT channel_id, channel_name FROM reports WHERE user_id = ? AND status = ?",
		userID,
		"open",
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var channelID, channelName string
	for rows.Next() {
		err = rows.Scan(&channelID, &channelName)
		if err != nil {
			panic(err)
		}
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}

	return channelID, channelName
}

func GetReportByChannelID(channelID string) (string, string) {
	rows, err := DBClient.Query(
		"SELECT user_id, channel_name FROM reports WHERE channel_id = ? AND status = ?",
		channelID,
		"open",
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var userID, channelName string
	for rows.Next() {
		err = rows.Scan(&userID, &channelName)
		if err != nil {
			panic(err)
		}
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}

	return userID, channelName
}

func checkUniqueOpenChannelName(channelName string) bool {
	rows, err := DBClient.Query(
		"SELECT channel_name FROM reports WHERE channel_name = ? AND status = ?",
		channelName,
		"open",
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var name string
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			panic(err)
		}
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}

	return name == ""
}

func getUniqueChannelName(channelName string) string {
	if !checkUniqueOpenChannelName(channelName) {
		// Get a number that is always 5 characters long
		randomNumber := 10000 + rand.Intn(90000) //nolint:gosec  // Only used for channel name generation
		channelName = channelName + " " + strconv.Itoa(randomNumber)
		return getUniqueChannelName(channelName)
	}

	return channelName
}

func CreateReport(channelID, channelName, userID string) {
	channelName = getUniqueChannelName(channelName)

	_, err := DBClient.Exec(
		`INSERT INTO reports (channel_id, channel_name, user_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		channelID,
		channelName,
		userID,
		"open",
		"now",
		"now",
	)
	if err != nil {
		panic(err)
	}
}

func CloseReport(channelID string) {
	_, err := DBClient.Exec(
		`UPDATE reports SET status = ?, updated_at = ? WHERE channel_id = ?`, "closed", "now",
		channelID,
	)
	if err != nil {
		panic(err)
	}
}
