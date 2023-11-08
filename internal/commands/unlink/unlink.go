package unlink

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slices"
	"ogusers-bot/internal/config"
	"ogusers-bot/internal/ogusers"
	"ogusers-bot/pkg/mongodb"
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
						Title:       "Unlink",
						Description: "You do not have permission to unlink users.",
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

	if len(options) < 1 {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Unlink",
					Description: "Please specify a uid/username or discord account to unlink.",
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})
		return
	}

	switch options[0].Name {
	case "uid":
		reply, err := mongodb.Database.FindOne(bson.M{"onsite": options[0].Value.(string), "verified": true}, mongodb.Database.Collections.Members)
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Unlink",
						Description: fmt.Sprintf("`%s` is not linked to a discord account.", options[0].Value.(string)),
						Color:       0x2B2D31,

						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
						},
					},
				},
			})
			return
		}

		s.GuildMemberRoleRemove(i.GuildID, reply["discord"].(string), config.C.Roles.Verified)
		s.GuildMemberRoleRemove(i.GuildID, reply["discord"].(string), config.C.Roles.Luminary)
		s.GuildMemberRoleRemove(i.GuildID, reply["discord"].(string), config.C.Roles.Revolution)

		mongodb.Database.DeleteOne(bson.M{"onsite": options[0].Value.(string), "verified": true}, mongodb.Database.Collections.Members)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Unlink",
					Description: fmt.Sprintf("Successfully unlinked `%s` from <@%s> (`%s`).", options[0].Value.(string), reply["discord"], reply["discord"]),
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})
	case "username":
		response, err := ogusers.GetUserByUsername(options[0].Value.(string))
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Unlink",
						Description: "Failed to get profile information.",
						Color:       0x2B2D31,

						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
						},
					},
				},
			})
			return
		} else if response.Error {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Unlink",
						Description: response.Message,
						Color:       0x2B2D31,

						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
						},
					},
				},
			})
			return
		}

		reply, err := mongodb.Database.FindOne(bson.M{"onsite": response.Uid, "verified": true}, mongodb.Database.Collections.Members)
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Unlink",
						Description: fmt.Sprintf("`%s` is not linked to a discord account.", options[0].Value.(string)),
						Color:       0x2B2D31,

						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
						},
					},
				},
			})
			return
		}

		s.GuildMemberRoleRemove(i.GuildID, reply["discord"].(string), config.C.Roles.Verified)
		s.GuildMemberRoleRemove(i.GuildID, reply["discord"].(string), config.C.Roles.Luminary)
		s.GuildMemberRoleRemove(i.GuildID, reply["discord"].(string), config.C.Roles.Revolution)

		mongodb.Database.DeleteOne(bson.M{"onsite": options[0].Value.(string), "verified": true}, mongodb.Database.Collections.Members)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Unlink",
					Description: fmt.Sprintf("Successfully unlinked `%s` from <@%s> (`%s`).", options[0].Value.(string), reply["discord"], reply["discord"]),
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})
	case "discord":
		reply, err := mongodb.Database.FindOne(bson.M{"discord": options[0].Value.(string), "verified": true}, mongodb.Database.Collections.Members)
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Unlink",
						Description: fmt.Sprintf("`%s` is not linked to a discord account.", options[0].Value.(string)),
						Color:       0x2B2D31,

						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
						},
					},
				},
			})
			return
		}

		s.GuildMemberRoleRemove(i.GuildID, reply["discord"].(string), config.C.Roles.Verified)
		s.GuildMemberRoleRemove(i.GuildID, reply["discord"].(string), config.C.Roles.Luminary)
		s.GuildMemberRoleRemove(i.GuildID, reply["discord"].(string), config.C.Roles.Revolution)

		mongodb.Database.DeleteOne(bson.M{"discord": options[0].Value.(string), "verified": true}, mongodb.Database.Collections.Members)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Unlink",
					Description: fmt.Sprintf("Successfully unlinked UID **`%s`** from <@%s> (`%s`).", reply["onsite"].(string), reply["discord"], reply["discord"]),
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})
	}
}
