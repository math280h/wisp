package moderation

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func WarnCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	Warn(i.Member.User.ID, s)

	// Respond to the command
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User has been warned",
		},
	})
	if err != nil {
		panic(err)
	}
}

func StrikeCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	Strike(i.Member.User.ID, s)

	// Respond to the command
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User has been striked",
		},
	})
	if err != nil {
		panic(err)
	}
}

func InfoCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	nick, points, reports := GetUserInfo(i.Member.User.ID)

	// Respond to command with embed
	embed := &discordgo.MessageEmbed{
		Color: 0x00ff00,
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
