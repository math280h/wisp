package db

type InfractionData struct {
	UserID      string
	Reason      string
	ModeratorID string
	CreatedAt   string
}

func GetAllInfractionsByUserID(user_id string) []InfractionData {
	// Return a list of all infractions for a user
	infractions := []InfractionData{}
	rows, err := DBClient.Query("SELECT user_id, reason, moderator_id, created_at FROM warns WHERE user_id = ?", user_id)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var infraction InfractionData
		err = rows.Scan(&infraction.UserID, &infraction.Reason, &infraction.ModeratorID, &infraction.CreatedAt)
		if err != nil {
			panic(err)
		}
		infractions = append(infractions, infraction)
	}

	return infractions
}

func GetMostRecentInfractionByUserID(user_id string) InfractionData {
	// Return the most recent infraction for a user
	var infraction InfractionData
	err := DBClient.QueryRow("SELECT user_id, reason, moderator_id, created_at FROM warns WHERE user_id = ? ORDER BY created_at DESC LIMIT 1", user_id).Scan(&infraction.UserID, &infraction.Reason, &infraction.ModeratorID, &infraction.CreatedAt)
	if err != nil {
		panic(err)
	}

	return infraction
}
