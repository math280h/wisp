package moderation

import "math280h/wisp/internal/db"

func GetUserInfo(userID string) (string, int, int) {
	// Get users name, warning points, and number of reports
	rows, err := db.DBClient.Query("SELECT nickname, points FROM users WHERE user_id = ?", userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var nickname string
	var points int
	for rows.Next() {
		err = rows.Scan(&nickname, &points)
		if err != nil {
			panic(err)
		}
	}

	// Get number of reports
	rows, err = db.DBClient.Query("SELECT COUNT(*) FROM reports WHERE user_id = ?", userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var reports int
	for rows.Next() {
		err = rows.Scan(&reports)
		if err != nil {
			panic(err)
		}
	}

	return nickname, points, reports
}
