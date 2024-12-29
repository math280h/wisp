package moderation

import (
	"context"
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
				// discordgo.Button{
				// 	Label:    "Notes",
				// 	Style:    discordgo.PrimaryButton,
				// 	CustomID: "notes" + suffix,
				// 	Emoji: &discordgo.ComponentEmoji{
				// 		Name: "📝",
				// 	},
				// },
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
		if err != db.ErrNotFound {
			log.Error().Err(err).Msg("Failed to get most recent infraction")
		}
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
		log.Error().Err(err).Msg("Failed to get user")
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
