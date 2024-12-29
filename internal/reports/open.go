package reports

import (
	"context"
	"errors"
	"fmt"
	"math280h/wisp/db"
	"math280h/wisp/internal/shared"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

const errMsg = "Error sending message to channel "

func OpenReportChannel(
	s *discordgo.Session,
	channelName string,
	originChannelID string,
	content string,
	authorID string,
	authorUsername string,
	avatarURL string,
	detain bool,
) (string, error) {
	// Create a new channel in the specified category
	channelData := &discordgo.GuildChannelCreateData{
		Name:     channelName,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: *shared.ReportCategory,
	}

	newChannel, err := s.GuildChannelCreateComplex(*shared.GuildID, *channelData)
	if err != nil {
		log.Error().Err(err).Msg("Error creating channel")
		return "", err
	}
	log.Info().Msg("Channel created ID:" + newChannel.ID + " Name:" + newChannel.Name + " User:" + authorID)
	// db.CreateReport(newChannel.ID, newChannel.Name, authorID)
	_, err = shared.DBClient.Report.CreateOne(
		db.Report.ChannelID.Set(newChannel.ID),
		db.Report.ChannelName.Set(newChannel.Name),
		db.Report.User.Link(
			db.User.UserID.Equals(authorID),
		),
	).Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to create report")
		return "", err
	}

	// Send embed message to the channel
	embed := &discordgo.MessageEmbed{
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User (Tag)",
				Value:  "<@" + authorID + ">",
				Inline: true,
			},
			{
				Name:   "User (Username)",
				Value:  authorUsername,
				Inline: true,
			},
			{
				Name:   "User (ID)",
				Value:  authorID,
				Inline: false,
			},
		},
		// Set the image as the users avatar
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: avatarURL,
		},
	}

	// Send an embed to the user saying that the channel has been created
	userEmbed := &discordgo.MessageEmbed{
		Color:       0x00ff00,
		Title:       "New report opened",
		Description: "A new report has been opened for you. Please use this channel to communicate with the staff.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Channel",
				Value:  "<#" + newChannel.ID + ">",
				Inline: true,
			},
		},
	}

	_, err = s.ChannelMessageSendEmbed(originChannelID, userEmbed)
	if err != nil {
		log.Error().Err(err).Msg(errMsg + originChannelID)
		return "", err
	}

	_, err = s.ChannelMessageSendEmbed(newChannel.ID, embed)
	if err != nil {
		log.Error().Err(err).Msg(errMsg + newChannel.ID)
		return "", err
	}

	// If detain flag is set to true, send embed saying that the user has been detained
	log.Debug().Msg("Detain flag: " + strconv.FormatBool(detain))
	if detain {
		embed = &discordgo.MessageEmbed{
			Color:       shared.Red,
			Title:       "Detained",
			Description: "You have been detained for the following reason: " + content,
		}
		_, err = s.ChannelMessageSendEmbed(newChannel.ID, embed)
		if err != nil {
			log.Error().Err(err).Msg(errMsg + newChannel.ID)
		}
	} else {
		// Send the users message to the channel
		_, err = s.ChannelMessageSend(newChannel.ID, content)
		if err != nil {
			log.Error().Err(err).Msg("Error sending message" + newChannel.ID)
			return "", err
		}
	}

	return newChannel.ID, nil
}

func OpenReport( //nolint:gocognit // This function is required to handle the message
	s *discordgo.Session,
	channelID string,
	authorID string,
	autherUsername string,
	content string,
	avatarURL string,
	detain bool,
) {
	// Remove any special characters from the channel name
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	expectedChannelName := re.ReplaceAllString(autherUsername, "")

	userObj, err := shared.DBClient.User.UpsertOne(
		db.User.UserID.Equals(authorID),
	).Create(
		db.User.UserID.Set(authorID),
		db.User.Nickname.Set(autherUsername),
	).Update().Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to open report")
		return
	}

	channelExists := true
	reportObj, err := shared.DBClient.Report.FindFirst(
		db.Report.UserID.Equals(userObj.ID),
		db.Report.Status.Equals("open"),
	).Exec(context.Background())
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			log.Error().Err(err).Msg("Failed to get report")
			return
		}

		channelExists = false
	}

	if !channelExists { //nolint:nestif // This is required to know how to handle the message
		// Create a new channel in the specified category
		_, reportErr := OpenReportChannel(
			s,
			expectedChannelName,
			channelID,
			content,
			authorID,
			autherUsername,
			avatarURL,
			detain,
		)
		if reportErr != nil {
			_, err = s.ChannelMessageSend(channelID, "There was an error sending your message. Please try again.")
			if err != nil {
				log.Error().Err(err).Msg(errMsg + channelID)
			}
		}
	} else {
		log.Debug().Msg("Open channel found, sending message to channel")

		if detain {
			embed := &discordgo.MessageEmbed{
				Color:       shared.Red,
				Title:       "Detained",
				Description: "You have been detained for the following reason: " + content,
			}
			_, chnlErr := s.ChannelMessageSendEmbed(reportObj.ChannelID, embed)
			if chnlErr != nil {
				log.Error().Err(err).Msg(errMsg + channelID)
			}
		} else {
			// Send the users message to the channel
			_, chnlErr := s.ChannelMessageSend(reportObj.ChannelID, content)
			if chnlErr != nil {
				log.Error().Err(err).Msg(errMsg + channelID)

				// If the error message contains Unknown Channel, attempt to create a new channel
				if fmt.Sprint(err) == "HTTP 404 Not Found, {\"message\": \"Unknown Channel\", \"code\": 10003}" {
					log.Debug().Msg("Channel not found, attempting to creating new channel")

					_, err = OpenReportChannel(
						s,
						expectedChannelName,
						channelID,
						content,
						authorID,
						autherUsername,
						avatarURL,
						detain,
					)
					if err != nil {
						_, err = s.ChannelMessageSend(channelID, "There was an error sending your message. Please try again.")
						if err != nil {
							log.Error().Err(err).Msg("Error sending message to channel " + channelID)
						}
					}
				}
				return
			}
		}
	}
}
