package db

import (
	"database/sql"
	"log"
	"math/rand"
	"strconv"
)

var db *sql.DB

func InitDb() {
	new_db, err := sql.Open("sqlite3", "./wisp.db")
	if err != nil {
		log.Fatal(err)
	}
	db = new_db

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS reports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		channel_id TEXT,
		channel_name TEXT,
		user_id TEXT,
		status TEXT,
		created_at TEXT,
		updated_at TEXT
	)`)
	if err != nil {
		panic(err)
	}
}

func GetReportByUserID(userID string) (string, string) {
	rows, err := db.Query("SELECT channel_id, channel_name FROM reports WHERE user_id = ? AND status = ?", userID, "open")
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

	return channelID, channelName
}

func GetReportByChannelID(channelID string) (string, string) {
	rows, err := db.Query("SELECT user_id, channel_name FROM reports WHERE channel_id = ? AND status = ?", channelID, "open")
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

	return userID, channelName
}

func checkUniqueOpenChannelName(channelName string) bool {
	rows, err := db.Query("SELECT channel_name FROM reports WHERE channel_name = ? AND status = ?", channelName, "open")
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

	return name == ""
}

func getUniqueChannelName(channelName string) string {
	if !checkUniqueOpenChannelName(channelName) {
		// Get a number that is always 5 characters long
		randomNumber := 10000 + rand.Intn(90000)
		channelName = channelName + " " + strconv.Itoa(randomNumber)
		return getUniqueChannelName(channelName)
	}

	return channelName
}

func CreateReport(channelID, channelName, userID string) {
	channelName = getUniqueChannelName(channelName)

	_, err := db.Exec(`INSERT INTO reports (channel_id, channel_name, user_id, status, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?)`, channelID, channelName, userID, "open", "now", "now")
	if err != nil {
		panic(err)
	}
}

func CloseReport(channelID string) {
	_, err := db.Exec(`UPDATE reports SET status = ?, updated_at = ? WHERE channel_id = ?`, "closed", "now", channelID)
	if err != nil {
		panic(err)
	}
}
