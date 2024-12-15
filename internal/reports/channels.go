package reports

import (
	"math280h/wisp/internal/db"
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func CreateReportChannel(expectedChannelName string, s *discordgo.Session, m *discordgo.MessageCreate) (channel_id string, err error) {
	// Create a new channel in the specified category
	channelData := &discordgo.GuildChannelCreateData{
		Name:     expectedChannelName,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: *shared.ReportCategory,
	}

	newChannel, err := s.GuildChannelCreateComplex(*shared.GuildID, *channelData)
	if err != nil {
		log.Error().Err(err).Msg("Error creating channel")
		return "", err
	}
	log.Info().Msg("Channel created ID:" + newChannel.ID + " Name:" + newChannel.Name + " User:" + m.Author.ID)
	db.CreateReport(newChannel.ID, newChannel.Name, m.Author.ID)

	// Send embed message to the channel
	embed := &discordgo.MessageEmbed{
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User (Tag)",
				Value:  "<@" + m.Author.ID + ">",
				Inline: true,
			},
			{
				Name:   "User (Username)",
				Value:  m.Author.Username,
				Inline: true,
			},
			{
				Name:   "User (ID)",
				Value:  m.Author.ID,
				Inline: false,
			},
		},
		// Set the image as the users avatar
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: m.Author.AvatarURL("256x256"),
		},
	}

	// Send an embed to the user saying that the channel has been created
	user_embed := &discordgo.MessageEmbed{
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
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, user_embed)
	if err != nil {
		log.Error().Err(err).Msg("Error sending message to channel " + m.ChannelID)
		return "", err
	}

	_, err = s.ChannelMessageSendEmbed(newChannel.ID, embed)
	if err != nil {
		log.Error().Err(err).Msg("Error sending message to channel " + newChannel.ID)
		return "", err
	}

	// Send the users message to the channel
	_, err = s.ChannelMessageSend(newChannel.ID, m.Content)
	if err != nil {
		log.Error().Err(err).Msg("Error sending message" + newChannel.ID)
		return "", err
	}

	return newChannel.ID, nil
}
