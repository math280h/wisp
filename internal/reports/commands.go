package reports

import (
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{ //nolint:gochecknoglobals // This is a list of commands for Discord
	{
		Name:                     "close",
		Description:              "Close the report",
		DefaultMemberPermissions: shared.Int64Ptr(discordgo.PermissionKickMembers),
	},
	{
		Name:                     "archive",
		Description:              "Archive the report",
		DefaultMemberPermissions: shared.Int64Ptr(discordgo.PermissionKickMembers),
	},
	{
		Name:                     "detain",
		Description:              "Detain the user",
		DefaultMemberPermissions: shared.Int64Ptr(discordgo.PermissionKickMembers),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to detain",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "The reason for the detainment",
				Required:    true,
			},
		},
	},
	{
		Name:                     "release",
		Description:              "Release the user",
		DefaultMemberPermissions: shared.Int64Ptr(discordgo.PermissionKickMembers),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to release",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "The reason for the release",
				Required:    true,
			},
		},
	},
}

var CommandHandlers = map[string]func( //nolint:gochecknoglobals // This is a map of commands to their handlers
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
){
	"close":   Close,
	"archive": Archive,
}
