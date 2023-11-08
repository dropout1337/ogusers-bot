package tickets

import (
	"github.com/bwmarrin/discordgo"
	"ogusers-bot/internal/discord"
)

var (
	allowInDirectMessages = false

	Command = discord.Command{
		Callable: handler,

		ApplicationCommand: &discordgo.ApplicationCommand{
			Name:        "tickets",
			Description: "Ticket commands",

			DMPermission: &allowInDirectMessages,

			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "send",
					Description: "Send the tickets embed in the specified channel",

					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionChannel,
							Name:        "channel",
							Description: "The channel to send the embed in",

							Required: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "add",
					Description: "Add a user to the ticket.",

					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionUser,
							Name:        "user",
							Description: "The user to add to the ticket.",

							Required: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "remove",
					Description: "Remove a user from the ticket.",

					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionUser,
							Name:        "user",
							Description: "The user to remove from the ticket.",

							Required: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "close",
					Description: "Close the ticket.",
				},
			},
		},
	}
	Components = []discord.Component{
		{
			CustomID: "general-support-ticket",
			Callable: createTicket("general-support-ticket"),
		},
		{
			CustomID: "account-support-ticket",
			Callable: createTicket("account-support-ticket"),
		},
		{
			CustomID: "middleman-request-ticket",
			Callable: createTicket("middleman-request-ticket"),
		},
		{
			CustomID: "purchase-ticket",
			Callable: createTicket("purchase-ticket"),
		},
		{
			CustomID: "close-ticket",
			Callable: closeTicket,
		},
	}
)
