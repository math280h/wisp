package db

func CreateUserIfNotExist(userID string, nickname string) {
	_, err := DBClient.Exec("INSERT OR IGNORE INTO users (user_id, nickname) VALUES (?, ?)", userID, nickname)
	if err != nil {
		panic(err)
	}
}
