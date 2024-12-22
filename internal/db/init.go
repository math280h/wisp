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
		user_id TEXT UNIQUE,
		nickname TEXT,
		points INTEGER DEFAULT 0
	)`)
	if err != nil {
		panic(err)
	}

	_, err = DBClient.Exec(`CREATE TABLE IF NOT EXISTS warns (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT,
		reason TEXT,
		moderator_id TEXT,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		panic(err)
	}

	_, err = DBClient.Exec(`CREATE TABLE IF NOT EXISTS suggestions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		embed_id TEXT,
		suggestion TEXT,
		user_id TEXT,
		status TEXT NOT NULL CHECK(status IN ('pending', 'approved', 'denied')),
		created_at TEXT DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		panic(err)
	}

	_, err = DBClient.Exec(`CREATE TABLE IF NOT EXISTS suggestion_votes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		suggestion_id INTEGER,
		user_id TEXT,
		sentiment TEXT NOT NULL CHECK(sentiment IN ('upvote', 'downvote')),
		created_at TEXT DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		panic(err)
	}
}
