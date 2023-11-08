package verify

import (
	"github.com/bwmarrin/discordgo"
	"ogusers-bot/internal/discord"
)

var (
	allowInDirectMessages = false

	Command = discord.Command{
		Callable: handler,

		ApplicationCommand: &discordgo.ApplicationCommand{
			Name:        "verify",
			Description: "Verify using your OGUsers account.",

			DMPermission: &allowInDirectMessages,

			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "Your OGUsers username.",

					Required: true,
				},
			},
		},
	}
)
