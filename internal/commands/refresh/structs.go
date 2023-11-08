package refresh

import (
	"github.com/bwmarrin/discordgo"
	"ogusers-bot/internal/discord"
)

var (
	allowInDirectMessages = false

	Command = discord.Command{
		Callable: handler,

		ApplicationCommand: &discordgo.ApplicationCommand{
			Name:        "refresh",
			Description: "Refresh your roles.",

			DMPermission: &allowInDirectMessages,
		},
	}
)
