package suggestions

import (
	"math280h/wisp/internal/shared"

	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{ //nolint:gochecknoglobals // This is a list of commands for Discord
	{
		Name:                     "suggest",
		Description:              "Suggest something for the server",
		DefaultMemberPermissions: shared.Int64Ptr(discordgo.PermissionViewChannel),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "suggestion",
				Description: "The suggestion",
				Required:    true,
			},
		},
	},
	{
		Name:                     "suggestion-status",
		Description:              "Set the status of a suggestion",
		DefaultMemberPermissions: shared.Int64Ptr(discordgo.PermissionAdministrator),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "suggestion_id",
				Description: "The ID of the suggestion",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "status",
				Description: "The status of the suggestion",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Approved",
						Value: "approved",
					},
					{
						Name:  "Denied",
						Value: "denied",
					},
				},
			},
		},
	},
}

var CommandHandlers = map[string]func( //nolint:gochecknoglobals // This is a map of commands to their handlers
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
){
	"suggest":           CreateSuggestionCommand,
	"suggestion-status": SetSuggestionStatusCommand,
}
