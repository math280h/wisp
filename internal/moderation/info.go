package moderation

import (
	"math280h/wisp/internal/db"
	"math280h/wisp/internal/shared"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func InfoCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	nick, points, reports := GetUserInfo(i.Member.User.ID)

	// Respond to command with embed
	embed := &discordgo.MessageEmbed{
		Color: shared.DarkBlue,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User (Tag)",
				Value:  "<@" + i.Member.User.ID + ">",
				Inline: true,
			},
			{
				Name:   "User (Username)",
				Value:  nick,
				Inline: true,
			},
			{
				Name:   "User (ID)",
				Value:  i.Member.User.ID,
				Inline: false,
			},
			{
				Name:   "Points",
				Value:  strconv.Itoa(points),
				Inline: true,
			},
			{
				Name:   "Reports",
				Value:  strconv.Itoa(reports),
				Inline: true,
			},
		},
		// Set the image as the users avatar
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: i.Member.User.AvatarURL("256x256"),
		},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		panic(err)
	}
}

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
