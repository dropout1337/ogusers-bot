package tags

import (
	"github.com/bwmarrin/discordgo"
	"ogusers-bot/internal/discord"
)

var (
	allowInDirectMessages = false

	Command = discord.Command{
		Callable: handler,

		ApplicationCommand: &discordgo.ApplicationCommand{
			Name:        "tags",
			Description: "View/create/delete tags.",

			DMPermission: &allowInDirectMessages,

			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionSubCommand,

					Name:        "list",
					Description: "View all tags.",
				},
				{
					Type: discordgo.ApplicationCommandOptionSubCommand,

					Name:        "view",
					Description: "View a tag.",

					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Tag name.",

							Required: true,
						},
					},
				},
				{
					Type: discordgo.ApplicationCommandOptionSubCommand,

					Name:        "create",
					Description: "Create a tag.",

					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Tag name.",

							Required: true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "content",
							Description: "Tag content.",

							Required: true,
						},
					},
				},
				{
					Type: discordgo.ApplicationCommandOptionSubCommand,

					Name:        "delete",
					Description: "Delete a tag.",

					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Tag name.",

							Required: true,
						},
					},
				},
			},
		},
	}
)
