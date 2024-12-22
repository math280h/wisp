package db

func CreateSuggestion(user_id string, suggestion string) int {
	suggestion_entry, err := DBClient.Exec("INSERT INTO suggestions (user_id, suggestion, status, embed_id) VALUES (?, ?, 'pending', '')", user_id, suggestion)
	if err != nil {
		panic(err)
	}

	// Return the ID of the result
	id, err := suggestion_entry.LastInsertId()
	if err != nil {
		panic(err)
	}

	return int(id)
}

func SetSuggestionEmbedID(suggestion_id int, embed_id string) {
	_, err := DBClient.Exec("UPDATE suggestions SET embed_id = ? WHERE id = ?", embed_id, suggestion_id)
	if err != nil {
		panic(err)
	}
}

func GetSuggestionByID(suggestion_id int) (string, string, string) {
	row := DBClient.QueryRow("SELECT embed_id, user_id, suggestion FROM suggestions WHERE id = ?", suggestion_id)

	var embed_id string
	var user_id string
	var suggestion string

	err := row.Scan(&embed_id, &user_id, &suggestion)
	if err != nil {
		return "", "", ""
	}

	return embed_id, user_id, suggestion
}

func SetSuggestionStatus(suggestion_id int, status string) {
	_, err := DBClient.Exec("UPDATE suggestions SET status = ? WHERE id = ?", status, suggestion_id)
	if err != nil {
		panic(err)
	}
}

func CreateSuggestionVote(suggestion_id int, user_id string, sentiment string) {
	_, err := DBClient.Exec("INSERT INTO suggestion_votes (suggestion_id, user_id, sentiment) VALUES (?, ?, ?)", suggestion_id, user_id, sentiment)
	if err != nil {
		panic(err)
	}
}

func GetSuggestionVoteByUserAndSuggestion(user_id string, suggestion_id int) (int, string) {
	row := DBClient.QueryRow("SELECT id, sentiment FROM suggestion_votes WHERE user_id = ? AND suggestion_id = ?", user_id, suggestion_id)

	var vote_id int
	var sentiment string

	err := row.Scan(&vote_id, &sentiment)
	if err != nil {
		return 0, ""
	}

	return vote_id, sentiment
}

func DeleteSuggestionVote(vote_id int) {
	_, err := DBClient.Exec("DELETE FROM suggestion_votes WHERE id = ?", vote_id)
	if err != nil {
		panic(err)
	}
}

func DeleteSuggestion(suggestion_id int) {
	_, err := DBClient.Exec("DELETE FROM suggestions WHERE id = ?", suggestion_id)
	if err != nil {
		panic(err)
	}
}

func GetSuggestionVoteCount(suggestion_id int) (int, int) {
	var upvotes int
	var downvotes int

	row := DBClient.QueryRow("SELECT COUNT(*) FROM suggestion_votes WHERE suggestion_id = ? AND sentiment = 'upvote'", suggestion_id)
	err := row.Scan(&upvotes)
	if err != nil {
		panic(err)
	}

	row = DBClient.QueryRow("SELECT COUNT(*) FROM suggestion_votes WHERE suggestion_id = ? AND sentiment = 'downvote'", suggestion_id)
	err = row.Scan(&downvotes)
	if err != nil {
		panic(err)
	}

	return upvotes, downvotes
}
