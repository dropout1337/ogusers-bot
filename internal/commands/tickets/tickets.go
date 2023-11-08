package tickets

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/slices"
	"ogusers-bot/internal/config"
	"ogusers-bot/pkg/logging"
	"strings"
	"time"
)

func handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	var authorized bool
	for _, role := range i.Interaction.Member.Roles {
		if slices.Contains(config.C.Roles.Staff, role) {
			authorized = true
		}
	}

	if !authorized {
		if i.Interaction.Member.Permissions&discordgo.PermissionAdministrator == 0 {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Tickets",
						Description: "You do not have the required permissions.",
						Color:       0x2B2D31,

						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
						},
					},
				},
			})
			return
		} else {
			authorized = true
		}
	}

	switch options[0].Name {
	case "send":
		var authorized bool
		for _, role := range i.Interaction.Member.Roles {
			if slices.Contains(config.C.Roles.Staff, role) {
				authorized = true
			}
		}

		if !authorized {
			if i.Interaction.Member.Permissions&discordgo.PermissionAdministrator == 0 {
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Embeds: &[]*discordgo.MessageEmbed{
						{
							Title:       "Tickets",
							Description: "You do not have the required permissions.",
							Color:       0x2B2D31,

							Thumbnail: &discordgo.MessageEmbedThumbnail{
								URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
							},
						},
					},
				})
				return
			} else {
				authorized = true
			}
		}

		channel, err := s.Channel(options[0].Options[0].Value.(string))
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Tickets",
						Description: "Invalid channel provided.",
						Color:       0x2B2D31,

						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
						},
					},
				},
			})
			return
		} else if channel.Type == discordgo.ChannelTypeGuildCategory {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Tickets",
						Description: "Invalid channel provided.",
						Color:       0x2B2D31,

						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
						},
					},
				},
			})
			return
		}

		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Tickets",
					Description: "Sending ticket embed in " + channel.Mention(),
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})

		s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:       "Ticket",
				Description: "Please press one of the buttons below that associate with your ticket.",
				Color:       0x2B2D31,

				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "General Support",
							Style:    discordgo.SecondaryButton,
							CustomID: "general-support-ticket",
						},
						discordgo.Button{
							Label:    "Account Support",
							Style:    discordgo.SecondaryButton,
							CustomID: "account-support-ticket",
						},
						discordgo.Button{
							Label:    "Middleman Request",
							Style:    discordgo.SecondaryButton,
							CustomID: "middleman-request-ticket",
						},
						discordgo.Button{
							Label:    "Purchase Ticket",
							Style:    discordgo.SecondaryButton,
							CustomID: "purchase-ticket",
						},
					},
				},
			},
		})
	case "add":
		var authorized bool
		for _, role := range i.Interaction.Member.Roles {
			if slices.Contains(config.C.Roles.Staff, role) {
				authorized = true
			}
		}

		if !authorized {
			if i.Interaction.Member.Permissions&discordgo.PermissionAdministrator == 0 {
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Embeds: &[]*discordgo.MessageEmbed{
						{
							Title:       "Tickets",
							Description: "You do not have the required permissions.",
							Color:       0x2B2D31,

							Thumbnail: &discordgo.MessageEmbedThumbnail{
								URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
							},
						},
					},
				})
				return
			} else {
				authorized = true
			}
		}

		s.ChannelPermissionSet(i.ChannelID, options[0].Options[0].Value.(string), discordgo.PermissionOverwriteTypeMember, discordgo.PermissionSendMessages|discordgo.PermissionViewChannel, 0)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Tickets",
					Description: "Added " + options[0].Options[0].Value.(string) + " to the ticket.",
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})
	case "remove":
		var authorized bool
		for _, role := range i.Interaction.Member.Roles {
			if slices.Contains(config.C.Roles.Staff, role) {
				authorized = true
			}
		}

		if !authorized {
			if i.Interaction.Member.Permissions&discordgo.PermissionAdministrator == 0 {
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Embeds: &[]*discordgo.MessageEmbed{
						{
							Title:       "Tickets",
							Description: "You do not have the required permissions.",
							Color:       0x2B2D31,

							Thumbnail: &discordgo.MessageEmbedThumbnail{
								URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
							},
						},
					},
				})
				return
			} else {
				authorized = true
			}
		}

		s.ChannelPermissionDelete(i.ChannelID, options[0].Options[0].Value.(string))
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Tickets",
					Description: "Removed " + options[0].Options[0].Value.(string) + " from the ticket.",
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})
	case "close":
		channel, _ := s.Channel(i.ChannelID)
		if channel.ParentID != config.C.Guild.TicketsCategory {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Tickets",
						Description: "This is not a ticket.",
						Color:       0x2B2D31,

						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
						},
					},
				},
			})
			return
		}

		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Tickets",
					Description: fmt.Sprintf("Closing ticket in <t:%v:R>..", time.Now().Add(10*time.Second).Unix()),
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})

		channel, _ = s.Channel(i.ChannelID)
		s.ChannelMessageSendComplex(config.C.Guild.TicketsLogsChannel, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: "Ticket",
				Color: 0x2B2D31,

				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Closed by",
						Value: fmt.Sprintf("%v#%v", i.Member.User.Username, i.Member.User.Discriminator),
					},
					{
						Name:  "Closed at",
						Value: fmt.Sprintf("<t:%v>", time.Now().Unix()),
					},
					{
						Name:  "Ticket",
						Value: fmt.Sprintf("#%v", channel.Name),
					},
				},

				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    fmt.Sprintf("Opened by %v#%v", i.Member.User.Username, i.Member.User.Discriminator),
							Style:    discordgo.SecondaryButton,
							CustomID: "idfk",
							Disabled: true,
						},
					},
				},
			},
		})

		time.Sleep(10 * time.Second)
		s.ChannelDelete(i.ChannelID)
	}
}

func createTicket(customID string) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})

		category, _ := s.Channel(config.C.Guild.TicketsCategory)
		channels, _ := s.GuildChannels(i.GuildID)

		for _, channel := range channels {
			if channel.ParentID != category.ID {
				continue
			}

			if strings.HasSuffix(channel.Name, strings.ReplaceAll(i.Member.User.Username, ".", "")) {
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Embeds: &[]*discordgo.MessageEmbed{
						{
							Title:       "Tickets",
							Description: "You already have a ticket open.",
							Color:       0x2B2D31,

							Thumbnail: &discordgo.MessageEmbedThumbnail{
								URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
							},
						},
					},
				})
				return
			}
		}

		var channelName string
		switch customID {
		case "general-support-ticket":
			channelName = "general-support-" + i.Member.User.Username
		case "account-support-ticket":
			channelName = "account-support-" + i.Member.User.Username
		case "middleman-request-ticket":
			channelName = "middleman-request-" + i.Member.User.Username
		case "purchase-ticket":
			channelName = "purchase-" + i.Member.User.Username
		default:
			channelName = "ticket-" + i.Member.User.Username
		}

		permissions := category.PermissionOverwrites
		permissions = append(permissions, &discordgo.PermissionOverwrite{
			ID:    i.Member.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionSendMessages | discordgo.PermissionViewChannel,
			Deny:  0x0,
		})

		permissions = append(permissions, &discordgo.PermissionOverwrite{
			ID:    "1115377891467866163",
			Type:  discordgo.PermissionOverwriteTypeRole,
			Allow: discordgo.PermissionSendMessages | discordgo.PermissionViewChannel,
			Deny:  0x0,
		})

		permissions = append(permissions, &discordgo.PermissionOverwrite{
			ID:    "1124146739851567205",
			Type:  discordgo.PermissionOverwriteTypeRole,
			Allow: discordgo.PermissionSendMessages | discordgo.PermissionViewChannel,
			Deny:  0x0,
		})

		channel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
			Name:                 channelName,
			Type:                 discordgo.ChannelTypeGuildText,
			ParentID:             config.C.Guild.TicketsCategory,
			PermissionOverwrites: permissions,
		})
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Tickets",
						Description: "Failed to create new ticket.",
						Color:       0x2B2D31,

						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
						},
					},
				},
			})
			return
		}

		if err != nil {
			logging.Logger.Error().
				Err(err).
				Str("guild", i.GuildID).
				Str("channel", channel.ID).
				Str("user", i.Member.User.ID).
				Msg("Failed to set channel permissions.")

			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Tickets",
						Description: "Failed to create new ticket.",
						Color:       0x2B2D31,

						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
						},
					},
				},
			})
			return
		}

		go func() {
			s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
				Content: i.Member.Mention(),
				Embed: &discordgo.MessageEmbed{
					Title:       "Ticket",
					Description: "Thank you for opening a ticket! Please wait for a staff member to assist you.",
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Close",
								Style:    discordgo.SecondaryButton,
								CustomID: "close-ticket",
							},
						},
					},
				},
			})
		}()

		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Tickets",
					Description: fmt.Sprintf("Opened ticket in <#%s>.", channel.ID),
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})

		s.ChannelMessageSendComplex(config.C.Guild.TicketsLogsChannel, &discordgo.MessageSend{
			Content: i.Member.Mention(),
			Embed: &discordgo.MessageEmbed{
				Title: "Ticket",
				Color: 0x2B2D31,

				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Opened at",
						Value: fmt.Sprintf("<t:%v>", time.Now().Unix()),
					},
					{
						Name:  "Ticket",
						Value: fmt.Sprintf("#%v", channel.Name),
					},
				},

				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    fmt.Sprintf("Opened by %v#%v", i.Member.User.Username, i.Member.User.Discriminator),
							Style:    discordgo.SecondaryButton,
							CustomID: "idfk",
							Disabled: true,
						},
					},
				},
			},
		})
	}
}

func closeTicket(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Tickets",
					Description: fmt.Sprintf("Closing ticket in <t:%v:R>..", time.Now().Add(10*time.Second).Unix()),
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		},
	})

	channel, _ := s.Channel(i.ChannelID)
	s.ChannelMessageSendComplex(config.C.Guild.TicketsLogsChannel, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Ticket",
			Color: 0x2B2D31,

			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Closed by",
					Value: fmt.Sprintf("%v#%v", i.Member.User.Username, i.Member.User.Discriminator),
				},
				{
					Name:  "Closed at",
					Value: fmt.Sprintf("<t:%v:R>", time.Now().Unix()),
				},
				{
					Name:  "Ticket",
					Value: fmt.Sprintf("#%v", channel.Name),
				},
			},

			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    fmt.Sprintf("Opened by %v#%v", i.Member.User.Username, i.Member.User.Discriminator),
						Style:    discordgo.SecondaryButton,
						CustomID: "idfk",
						Disabled: true,
					},
				},
			},
		},
	})

	time.Sleep(5 * time.Second)
	s.ChannelDelete(i.ChannelID)
}
