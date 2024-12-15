package db

import (
	"database/sql"
	"log"
)

var DBClient *sql.DB

func InitDb() {
	new_db, err := sql.Open("sqlite3", "./wisp.db")
	if err != nil {
		log.Fatal(err)
	}
	DBClient = new_db

	_, err = DBClient.Exec(`CREATE TABLE IF NOT EXISTS reports (
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

	_, err = DBClient.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT,
		nickname TEXT,
		points INTEGER DEFAULT 0
	)`)
	if err != nil {
		panic(err)
	}
}
