package discord

import "github.com/bwmarrin/discordgo"

type Discord struct {
	Session *discordgo.Session

	Commands          map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	MessageComponents map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)

	RegisteredCommands []*discordgo.ApplicationCommand
}

type Command struct {
	Callable           func(s *discordgo.Session, i *discordgo.InteractionCreate)
	ApplicationCommand *discordgo.ApplicationCommand
}

type Component struct {
	CustomID string
	Callable func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
