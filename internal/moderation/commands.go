package moderation

import (
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{ //nolint:gochecknoglobals // This is a list of commands for Discord
	{
		Name:                     "note",
		Description:              "Add a note to a user",
		DefaultMemberPermissions: shared.Int64Ptr(discordgo.PermissionKickMembers),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to add a note to",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "The text of the note",
				Required:    true,
			},
		},
	},
	{
		Name:                     "warn",
		Description:              "Warn the user",
		DefaultMemberPermissions: shared.Int64Ptr(discordgo.PermissionKickMembers),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to warn",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "The reason for the warning",
				Required:    true,
			},
		},
	},
	{
		Name:                     "strike",
		Description:              "Strike the user",
		DefaultMemberPermissions: shared.Int64Ptr(discordgo.PermissionKickMembers),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to strike",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "The reason for the strike",
				Required:    true,
			},
		},
	},
	{
		Name:                     "kick",
		Description:              "Kick the user",
		DefaultMemberPermissions: shared.Int64Ptr(discordgo.PermissionKickMembers),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to kick",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "The reason for the kick",
				Required:    true,
			},
		},
	},
	{
		Name:                     "ban",
		Description:              "Ban the user",
		DefaultMemberPermissions: shared.Int64Ptr(discordgo.PermissionBanMembers),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to ban",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "The reason for the ban",
				Required:    true,
			},
		},
	},
	{
		Name:                     "remove-infraction",
		Description:              "Remove an infraction from a user",
		DefaultMemberPermissions: shared.Int64Ptr(discordgo.PermissionAdministrator),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "infraction_id",
				Description: "The ID of the infraction to remove",
				Required:    true,
			},
		},
	},
	{
		Name:                     "info",
		Description:              "Get information about a user",
		DefaultMemberPermissions: shared.Int64Ptr(discordgo.PermissionKickMembers),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to get information about",
				Required:    true,
			},
		},
	},
}

var CommandHandlers = map[string]func( //nolint:gochecknoglobals // This is a map of commands to their handlers
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
){
	"note":              NoteCommand,
	"warn":              WarnCommand,
	"strike":            StrikeCommand,
	"info":              InfoCommand,
	"detain":            DetainUserCommand,
	"release":           ReleaseUserCommand,
	"remove-infraction": RemoveInfraction,
}
