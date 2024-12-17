package moderation

import (
	"math280h/wisp/internal/shared"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func generateInfractionEmbed(reason string, points int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       0xe74c3c,
		Title:       "You have been warned",
		Description: "You have been warned by a moderator. Please ensure you follow the rules in the future.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Guild",
				Value: shared.GuildName,
			},
			{
				Name:   "Reason",
				Value:  reason,
				Inline: true,
			},
			{
				Name:   "Total Points",
				Value:  strconv.Itoa(points) + " / " + strconv.Itoa(*shared.MaxPoints),
				Inline: true,
			},
		},
	}
}

func WarnCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	reason := i.ApplicationCommandData().Options[1].StringValue()
	points := Warn(i.Member.User.ID, reason, i.Member.User.ID, s)

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

	// Inform the user that they have been warned
	embed := generateInfractionEmbed(reason, points)

	userChannel, err := s.UserChannelCreate(i.Member.User.ID)
	if err != nil {
		panic(err)
	}

	_, err = s.ChannelMessageSendEmbed(userChannel.ID, embed)
	if err != nil {
		panic(err)
	}
}

func StrikeCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	reason := i.ApplicationCommandData().Options[1].StringValue()
	points := Strike(i.Member.User.ID, reason, i.Member.User.ID, s)

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

	// Inform the user that they have been striked
	embed := generateInfractionEmbed(reason, points)

	userChannel, err := s.UserChannelCreate(i.Member.User.ID)
	if err != nil {
		panic(err)
	}

	_, err = s.ChannelMessageSendEmbed(userChannel.ID, embed)
	if err != nil {
		panic(err)
	}
}

func KickCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	reason := i.ApplicationCommandData().Options[1].StringValue()
	Kick(i.Member.User.ID, reason, i.Member.User.ID, s)

	// Respond to the command
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User has been kicked",
		},
	})
	if err != nil {
		panic(err)
	}
}

func BanCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	reason := i.ApplicationCommandData().Options[1].StringValue()
	Ban(i.Member.User.ID, reason, i.Member.User.ID, s)

	// Respond to the command
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User has been banned",
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
