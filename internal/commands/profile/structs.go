package profile

import (
	"github.com/bwmarrin/discordgo"
	"ogusers-bot/internal/discord"
)

var (
	allowInDirectMessages = false

	Command = discord.Command{
		Callable: handler,

		ApplicationCommand: &discordgo.ApplicationCommand{
			Name:        "profile",
			Description: "View profile user profile information.",

			DMPermission: &allowInDirectMessages,

			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionSubCommand,

					Name:        "uid",
					Description: "View profile information by user ID.",

					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "uid",
							Description: "User ID to view profile information for.",

							Required: true,
						},
					},
				},
				{
					Type: discordgo.ApplicationCommandOptionSubCommand,

					Name:        "discord",
					Description: "View profile information by Discord account.",

					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionUser,
							Name:        "user",
							Description: "Discord user to view profile information for.",

							Required: true,
						},
					},
				},
				{
					Type: discordgo.ApplicationCommandOptionSubCommand,

					Name:        "username",
					Description: "View profile information by username.",

					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "username",
							Description: "Username to view profile information for.",

							Required: true,
						},
					},
				},
			},
		},
	}
)
