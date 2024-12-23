package db

func CreateSuggestion(userID string, suggestion string) int {
	suggestionEntry, err := DBClient.Exec(
		"INSERT INTO suggestions (user_id, suggestion, status, embed_id) VALUES (?, ?, 'pending', '')",
		userID,
		suggestion,
	)
	if err != nil {
		panic(err)
	}

	// Return the ID of the result
	id, err := suggestionEntry.LastInsertId()
	if err != nil {
		panic(err)
	}

	return int(id)
}

func SetSuggestionEmbedID(suggestionID int, embedID string) {
	_, err := DBClient.Exec("UPDATE suggestions SET embed_id = ? WHERE id = ?", embedID, suggestionID)
	if err != nil {
		panic(err)
	}
}

func GetSuggestionByID(suggestionID int) (string, string, string) {
	row := DBClient.QueryRow(
		"SELECT embed_id, user_id, suggestion FROM suggestions WHERE id = ?",
		suggestionID,
	)

	var embedID string
	var userID string
	var suggestion string

	err := row.Scan(&embedID, &userID, &suggestion)
	if err != nil {
		return "", "", ""
	}

	return embedID, userID, suggestion
}

func SetSuggestionStatus(suggestionID int, status string) {
	_, err := DBClient.Exec("UPDATE suggestions SET status = ? WHERE id = ?", status, suggestionID)
	if err != nil {
		panic(err)
	}
}

func CreateSuggestionVote(suggestionID int, userID string, sentiment string) {
	_, err := DBClient.Exec(
		"INSERT INTO suggestion_votes (suggestion_id, user_id, sentiment) VALUES (?, ?, ?)",
		suggestionID,
		userID,
		sentiment,
	)
	if err != nil {
		panic(err)
	}
}

func GetSuggestionVoteByUserAndSuggestion(userID string, suggestionID int) (int, string) {
	row := DBClient.QueryRow(
		"SELECT id, sentiment FROM suggestion_votes WHERE user_id = ? AND suggestion_id = ?",
		userID,
		suggestionID,
	)

	var voteID int
	var sentiment string

	err := row.Scan(&voteID, &sentiment)
	if err != nil {
		return 0, ""
	}

	return voteID, sentiment
}

func DeleteSuggestionVote(voteID int) {
	_, err := DBClient.Exec("DELETE FROM suggestion_votes WHERE id = ?", voteID)
	if err != nil {
		panic(err)
	}
}

func DeleteSuggestion(suggestionID int) {
	_, err := DBClient.Exec("DELETE FROM suggestions WHERE id = ?", suggestionID)
	if err != nil {
		panic(err)
	}
}

func GetSuggestionVoteCount(suggestionID int) (int, int) {
	var upvotes int
	var downvotes int

	row := DBClient.QueryRow(
		"SELECT COUNT(*) FROM suggestion_votes WHERE suggestion_id = ? AND sentiment = 'upvote'",
		suggestionID,
	)
	err := row.Scan(&upvotes)
	if err != nil {
		panic(err)
	}

	row = DBClient.QueryRow(
		"SELECT COUNT(*) FROM suggestion_votes WHERE suggestion_id = ? AND sentiment = 'downvote'",
		suggestionID,
	)
	err = row.Scan(&downvotes)
	if err != nil {
		panic(err)
	}

	return upvotes, downvotes
}
