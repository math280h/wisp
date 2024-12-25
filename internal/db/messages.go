package db

import "github.com/rs/zerolog/log"

func CreateMessage(id string, content string, authorID string, authorTag string, timestamp string, channelID string) {
	_, err := DBClient.Exec(
		"INSERT INTO messages (id, content, author_id, author_tag, timestamp, channel_id) VALUES (?, ?, ?, ?, ?, ?)",
		id,
		content,
		authorID,
		authorTag,
		timestamp,
		channelID,
	)
	if err != nil {
		log.Error().Msg("Failed to save message")
	}
}

func UpdateContentByID(id string, content string) {
	_, err := DBClient.Exec("UPDATE messages SET content = ? WHERE id = ?", content, id)
	if err != nil {
		log.Error().Msg("Failed to save message edit")
	}
}

func GetMessageByID(id string) (string, string, string, string, string, string) {
	row := DBClient.QueryRow(
		"SELECT content, author_id, author_tag, timestamp, channel_id, created_at FROM messages WHERE id = ?",
		id,
	)

	var content string
	var authorID string
	var authorTag string
	var timestamp string
	var channelID string
	var createdAt string

	err := row.Scan(&content, &authorID, &authorTag, &timestamp, &channelID, &createdAt)
	if err != nil {
		return "", "", "", "", "", ""
	}

	return content, authorID, authorTag, timestamp, channelID, createdAt
}

func DeleteMessageByID(id string) {
	_, err := DBClient.Exec("DELETE FROM messages WHERE id = ?", id)
	if err != nil {
		log.Error().Msg("Failed to delete message")
	}
}
