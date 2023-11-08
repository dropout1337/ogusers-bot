package commands

import (
	"ogusers-bot/internal/commands/fake"
	"ogusers-bot/internal/commands/forceverify"
	"ogusers-bot/internal/commands/profile"
	"ogusers-bot/internal/commands/refresh"
	"ogusers-bot/internal/commands/tags"
	"ogusers-bot/internal/commands/tickets"
	"ogusers-bot/internal/commands/unlink"
	"ogusers-bot/internal/commands/verify"
	"ogusers-bot/internal/discord"
	"ogusers-bot/pkg/logging"
)

var (
	commands = []discord.Command{
		profile.Command,
		verify.Command,
		refresh.Command,
		unlink.Command,
		forceverify.Command,

		tags.Command,
		tickets.Command,
		fake.Command,
	}

	components []discord.Component
)

func init() {
	components = append(components, tickets.Components...)
}

func Hook(c *discord.Discord) {
	for _, component := range components {
		err := c.AddHandler(component)
		if err != nil {
			logging.Logger.Fatal().Err(err).Msg("Failed to add handler")
		} else {
			logging.Logger.Info().Str("component", component.CustomID).Msg("Added component handler")
		}
	}

	for _, command := range commands {
		err := c.AddHandler(command)
		if err != nil {
			logging.Logger.Fatal().Err(err).Msg("Failed to add handler")
		} else {
			logging.Logger.Info().Str("command", command.ApplicationCommand.Name).Msg("Added command handler")
		}
	}
}
