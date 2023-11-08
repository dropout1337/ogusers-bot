package unlink

import (
	"github.com/bwmarrin/discordgo"
	"ogusers-bot/internal/discord"
)

var (
	allowInDirectMessages = false

	Command = discord.Command{
		Callable: handler,

		ApplicationCommand: &discordgo.ApplicationCommand{
			Name:        "unlink",
			Description: "Unlink a discord account from an OGUsers account",

			DMPermission: &allowInDirectMessages,

			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "OGUsers username.",

					Required: false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "uid",
					Description: "OGUsers uid.",

					Required: false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "discord",
					Description: "Discord account.",

					Required: false,
				},
			},
		},
	}
)
