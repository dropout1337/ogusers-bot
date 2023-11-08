package fake

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/slices"
	"ogusers-bot/internal/config"
	"strings"
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
						Title:       "Fake",
						Description: "You do not have permission to force fake messages.",
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

	var webhook *discordgo.Webhook

	// discord.go hasn't added global name support...
	body, err := s.RequestWithBucketID("GET", discordgo.EndpointGuildMember(i.GuildID, options[0].Value.(string)), nil, discordgo.EndpointGuildMember(i.GuildID, ""))
	if err != nil {
		return
	}
	payload := gjson.ParseBytes(body)

	user, _ := s.GuildMember(i.GuildID, options[0].Value.(string))
	webhooks, _ := s.ChannelWebhooks(i.ChannelID)

	if len(webhooks) == 0 {
		webhook, _ = s.WebhookCreate(i.ChannelID, "/fake", "")
	} else {
		webhook = webhooks[0]
	}

	var content string
	content = strings.ReplaceAll(options[1].Value.(string), "@everyone", "@\u200beveryone")
	content = strings.ReplaceAll(options[1].Value.(string), "@here", "@\u200bhere")
	content = strings.ReplaceAll(options[1].Value.(string), "<@&", "<@\u200b&")

	s.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
		Content:   content,
		Username:  payload.Get("user.display_name").String(),
		AvatarURL: user.User.AvatarURL(""),
	})

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "Fake",
				Description: "Successfully sent fake message.",
				Color:       0x2B2D31,

				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
				},
			},
		},
	})
}
