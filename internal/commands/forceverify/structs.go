package forceverify

import (
	"github.com/bwmarrin/discordgo"
	"ogusers-bot/internal/discord"
)

var (
	allowInDirectMessages = false

	Command = discord.Command{
		Callable: handler,

		ApplicationCommand: &discordgo.ApplicationCommand{
			Name:        "force-verify",
			Description: "Force-verify a discord account.",

			DMPermission: &allowInDirectMessages,

			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to force-verify.",

					Required: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "uid",
					Description: "Their OGUsers uid.",

					Required: true,
				},
			},
		},
	}
)
