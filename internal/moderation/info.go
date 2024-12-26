package moderation

import (
	"math280h/wisp/internal/db"
	"math280h/wisp/internal/shared"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func GenerateInfoButtons(channelID string, embedID string, userID string) []discordgo.MessageComponent {
	var suffix = ":" + userID + ":" + embedID + ":" + channelID
	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Overview",
					Style:    discordgo.PrimaryButton,
					CustomID: "overview" + suffix,
					Emoji: &discordgo.ComponentEmoji{
						Name: "🔍",
					},
				},
				discordgo.Button{
					Label:    "Infractions",
					Style:    discordgo.PrimaryButton,
					CustomID: "infractions" + suffix,
					Emoji: &discordgo.ComponentEmoji{
						Name: "⚠️",
					},
				},
				discordgo.Button{
					Label:    "Notes",
					Style:    discordgo.PrimaryButton,
					CustomID: "notes" + suffix,
					Emoji: &discordgo.ComponentEmoji{
						Name: "📝",
					},
				},
				discordgo.Button{
					Label:    "Messages",
					Style:    discordgo.PrimaryButton,
					CustomID: "messages" + suffix,
					Emoji: &discordgo.ComponentEmoji{
						Name: "💬",
					},
				},
				discordgo.Button{
					Label:    "Leaves",
					Style:    discordgo.PrimaryButton,
					CustomID: "leaves" + suffix,
					Emoji: &discordgo.ComponentEmoji{
						Name: "🚪",
					},
				},
			},
		},
	}
	return buttons
}

func GenerateOverviewEmbed(userID string, nickname string, points int, reports int, avatar string) discordgo.MessageEmbed {
	mostRecentInfraction := db.GetMostRecentInfractionByUserID(userID)

	embed := discordgo.MessageEmbed{
		Color: shared.DarkBlue,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User (Tag)",
				Value:  "<@" + userID + ">",
				Inline: true,
			},
			{
				Name:   "User (Username)",
				Value:  nickname,
				Inline: true,
			},
			{
				Name:   "User (ID)",
				Value:  userID,
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
			{
				Name:   "__Most Recent Infraction__",
				Value:  "",
				Inline: false,
			},
			{
				Name: "Staff ::" + "<@" + mostRecentInfraction.ModeratorID + ">",
				Value: "Type: **Strike** \n" +
					"Date: **" + strings.Split(mostRecentInfraction.CreatedAt, "T")[0] + "** (" + shared.StringTimeToDiscordTimestamp(mostRecentInfraction.CreatedAt) + ")\n" +
					"Reason: " + mostRecentInfraction.Reason,
				Inline: false,
			},
		},
		// Set the image as the users avatar
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: avatar,
		},
	}
	return embed
}

func InfoCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.ApplicationCommandData().Options[0].UserValue(s)

	nick, points, reports := GetUserInfo(user.ID)

	// Respond to command with embed
	embed := GenerateOverviewEmbed(user.ID, nick, points, reports, user.AvatarURL("256x256"))
	infoMessage, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Embed: &embed,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send info message")
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to get user info",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to send fail interaction response for info")
		}
		return
	}

	buttons := GenerateInfoButtons(i.ChannelID, infoMessage.ID, user.ID)
	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    i.ChannelID,
		ID:         infoMessage.ID,
		Components: &buttons,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to add buttons to info message")
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User information available below",
			Flags:   64,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send interaction response for info")
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
