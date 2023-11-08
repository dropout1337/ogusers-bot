package main

import (
	"github.com/bwmarrin/discordgo"
	"ogusers-bot/internal/commands"
	"ogusers-bot/internal/config"
	"ogusers-bot/internal/discord"
	"ogusers-bot/internal/events"
	"ogusers-bot/pkg/logging"
)

func main() {
	client, err := discord.New(config.C.Discord.Token)
	if err != nil {
		logging.Logger.Fatal().Err(err).Msg("Failed to create Discord client")
	}

	discord.C = client

	client.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Status: config.C.Discord.Status.Type,
			Activities: []*discordgo.Activity{
				{
					Name: config.C.Discord.Status.Text,
					Type: discordgo.ActivityTypeListening,
				},
			},
		})

		logging.Logger.Info().Str("username", s.State.User.Username).Msg("Logged in")
		events.ChannelPurger(s)
		client.Session.AddHandler(events.OnMessage)
		commands.Hook(client)
	})

	err = client.Open()
	if err != nil {
		logging.Logger.Fatal().Err(err).Msg("Failed to open Discord client")
	}
}
