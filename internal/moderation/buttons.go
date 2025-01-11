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

func InfoButtons(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	// Get the CustomID
	var customID = i.MessageComponentData().CustomID
	// Split it by : to get the values
	// The format is <action>:<userID>:<embedID>:<channelID>
	var customIDSplit = strings.Split(customID, ":")
	if len(customIDSplit) != 4 {
		return
	}

	// Get the action
	var action = customIDSplit[0]

	// Get the userID
	var userID = customIDSplit[1]
	// Get user from the user ID
	discordUser, err := s.User(userID)
	if err != nil {
		shared.SimpleEphemeralInteractionResponse("Failed to get discord user", s, i.Interaction)
		log.Error().Err(err).Msg("Failed to get discord user")
		return
	}

	// Get the embedID
	var embedID = customIDSplit[2]
	// Get the channelID
	var channelID = customIDSplit[3]

	userObj, err := shared.GetUserIfExists(&discordgo.User{
		ID:       userID,
		Username: discordUser.Username,
	})
	if err != nil {
		shared.SimpleEphemeralInteractionResponse("Failed to get or create user", s, i.Interaction)
		log.Error().Err(err).Msg("Failed to get user")
		return
	}

	reportCount := getUserReportCount(userObj.ID)

	log.Debug().Msg("Got user info button event")

	var embed = discordgo.MessageEmbed{
		Title: "Infractions",
		Color: shared.DarkBlue,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User (Tag)",
				Value:  "",
				Inline: true,
			},
			{
				Name:   "Points",
				Value:  strconv.Itoa(userObj.Points),
				Inline: true,
			},
			{
				Name:   "Reports",
				Value:  strconv.Itoa(reportCount),
				Inline: true,
			},
		},
		// Set the image as the users avatar
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: discordUser.AvatarURL("256x256"),
		},
	}

	switch action {
	case "overview":
		embed = GenerateOverviewEmbed(*userObj, userID, reportCount, discordUser.AvatarURL("256x256"))
	case "infractions":
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "__Infractions__",
			Value: "",
		})

		// Get all infractions
		infractions, infractionErr := shared.DBClient.Infraction.FindMany(
			db.Infraction.UserID.Equals(userObj.ID),
		).OrderBy(
			db.Infraction.CreatedAt.Order(db.SortOrderDesc),
		).Exec(context.Background())
		if infractionErr != nil {
			shared.SimpleEphemeralInteractionResponse("Failed to get infractions", s, i.Interaction)
			log.Error().Err(err).Msg("Failed to get infractions")
			return
		}

		for _, infraction := range infractions {
			dateWithoutTime := strings.Split(infraction.CreatedAt, "T")[0]

			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name: "ID:: " + strconv.Itoa(infraction.ID) + " :: Staff ::" + infraction.ModeratorUsername,
				Value: "Type: **" + infraction.Type + "** \n" +
					"Date: **" + dateWithoutTime + "** (" + shared.StringTimeToDiscordTimestamp(infraction.CreatedAt) + ")\n" +
					"Reason: " + infraction.Reason,
				Inline: false,
			})
		}
	case "notes":
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "__Notes__",
			Value: "",
		})

		// Get all notes
		notes, noteErr := shared.DBClient.Note.FindMany(
			db.Note.UserID.Equals(userObj.ID),
		).OrderBy(
			db.Note.CreatedAt.Order(db.SortOrderDesc),
		).Exec(context.Background())
		if noteErr != nil {
			shared.SimpleEphemeralInteractionResponse("Failed to get notes", s, i.Interaction)
			log.Error().Err(err).Msg("Failed to get notes")
			return
		}

		for _, note := range notes {
			dateWithoutTime := strings.Split(note.CreatedAt, "T")[0]

			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name: "ID:: " + strconv.Itoa(note.ID) + " :: Staff ::" + note.ModeratorUsername,
				Value: "Date: **" + dateWithoutTime + "** (" + shared.StringTimeToDiscordTimestamp(note.CreatedAt) + ")\n" +
					"Note: " + note.Content,
				Inline: false,
			})
		}
	case "messages":
		// Show the messages
	case "leaves":
		// Show the leaves
	default:
		return
	}

	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel: channelID,
		ID:      embedID,
		Embeds:  &[]*discordgo.MessageEmbed{&embed},
	})
	if err != nil {
		shared.SimpleEphemeralInteractionResponse("Failed to edit message", s, i.Interaction)
		log.Error().Err(err).Msg("Failed to edit message")
		return
	}

	// Respond to the interaction
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to respond to interaction")
		return
	}
}
