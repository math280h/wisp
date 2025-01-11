package moderation

import (
	"context"
	"errors"
	"math280h/wisp/db"
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
				// discordgo.Button{
				// 	Label:    "Messages",
				// 	Style:    discordgo.PrimaryButton,
				// 	CustomID: "messages" + suffix,
				// 	Emoji: &discordgo.ComponentEmoji{
				// 		Name: "💬",
				// 	},
				// },
				// discordgo.Button{
				// 	Label:    "Leaves",
				// 	Style:    discordgo.PrimaryButton,
				// 	CustomID: "leaves" + suffix,
				// 	Emoji: &discordgo.ComponentEmoji{
				// 		Name: "🚪",
				// 	},
				// },
			},
		},
	}
	return buttons
}

func GenerateOverviewEmbed(user db.UserModel, userDiscordID string, reports int, avatar string) discordgo.MessageEmbed {
	mostRecentInfraction, err := shared.DBClient.Infraction.FindFirst(
		db.Infraction.UserID.Equals(user.ID),
	).OrderBy(
		db.Infraction.CreatedAt.Order(db.SortOrderDesc),
	).Exec(context.Background())
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			log.Error().Err(err).Msg("Failed to get most recent infraction")
		}
	}

	// Ensure last leave is not empty, if it is set to "Never"
	usrLastLeave, lastLeaveOk := user.LastLeave()
	if !lastLeaveOk {
		usrLastLeave = "Never"
	}

	usrLastJoin, lastJoinOk := user.LastJoin()
	if !lastJoinOk {
		usrLastJoin = "Never"
	}

	embed := discordgo.MessageEmbed{
		Color: shared.DarkBlue,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User (Tag)",
				Value:  "<@" + userDiscordID + ">",
				Inline: true,
			},
			{
				Name:   "User (Username)",
				Value:  user.Nickname,
				Inline: true,
			},
			{
				Name:   "User (ID)",
				Value:  strconv.Itoa(user.ID),
				Inline: false,
			},
			{
				Name:   "Points",
				Value:  strconv.Itoa(user.Points),
				Inline: true,
			},
			{
				Name:   "Reports",
				Value:  strconv.Itoa(reports),
				Inline: true,
			},
			{
				Name:   "__Leave History__",
				Value:  "",
				Inline: false,
			},
			{
				Name:   "Last Join",
				Value:  shared.StringWithTzToDiscordTimestamp(usrLastJoin),
				Inline: true,
			},
			{
				Name:   "Last Leave",
				Value:  shared.StringWithTzToDiscordTimestamp(usrLastLeave),
				Inline: true,
			},
			{
				Name:   "Leave Count",
				Value:  strconv.Itoa(user.LeaveCount),
				Inline: true,
			},
			{
				Name:   "__Most Recent Infraction__",
				Value:  "",
				Inline: false,
			},
		},
		// Set the image as the users avatar
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: avatar,
		},
	}

	if mostRecentInfraction != nil {
		// Add the most recent infraction to the embed
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "Staff ::" + mostRecentInfraction.ModeratorUsername,
			Value: "Type: **" + mostRecentInfraction.Type + "** \n" +
				// 2024-12-29 05:10:23
				"Date: **" + strings.Split(mostRecentInfraction.CreatedAt, " ")[0] +
				"** (" + shared.StringTimeToDiscordTimestamp(mostRecentInfraction.CreatedAt) + ")\n" +
				"Reason: " + mostRecentInfraction.Reason,
			Inline: false,
		})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "No infractions",
			Value:  "This user has no infractions",
			Inline: false,
		})
	}

	return embed
}

func InfoCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.ApplicationCommandData().Options[0].UserValue(s)

	userObj, err := shared.DBClient.User.FindUnique(
		db.User.UserID.Equals(user.ID),
	).Exec(context.Background())
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			shared.SimpleEphemeralInteractionResponse("User not found", s, i.Interaction)
			return
		}
		shared.SimpleEphemeralInteractionResponse("Failed to get user", s, i.Interaction)
		log.Error().Err(err).Msg("Failed to get user")
		return
	}

	reportCount := getUserReportCount(userObj.ID)

	// Respond to command with embed
	embed := GenerateOverviewEmbed(
		*userObj,
		user.ID,
		reportCount,
		user.AvatarURL("256x256"),
	)
	infoMessage, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Embed: &embed,
	})
	if err != nil {
		shared.SimpleEphemeralInteractionResponse("Failed to send info message", s, i.Interaction)
		log.Error().Err(err).Msg("Failed to send info message")
		return
	}

	buttons := GenerateInfoButtons(i.ChannelID, infoMessage.ID, user.ID)
	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    i.ChannelID,
		ID:         infoMessage.ID,
		Components: &buttons,
	})
	if err != nil {
		shared.SimpleEphemeralInteractionResponse("Failed to add buttons to info message", s, i.Interaction)
		log.Error().Err(err).Msg("Failed to add buttons to info message")
		return
	}

	shared.SimpleEphemeralInteractionResponse("User info generated below", s, i.Interaction)
}
