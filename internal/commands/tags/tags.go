package tags

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slices"
	"ogusers-bot/internal/config"
	"ogusers-bot/pkg/mongodb"
	"strings"
)

func handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsLoading,
		},
	})

	switch options[0].Name {
	case "view":
		viewTag(s, i, options[0].Options[0].Value.(string))
	case "create":
		createTag(s, i, options[0].Options[0].Value.(string), options[0].Options[1].Value.(string))
	case "delete":
		deleteTag(s, i, options[0].Options[0].Value.(string))
	case "list":
		listTags(s, i)
	}
}

func viewTag(s *discordgo.Session, i *discordgo.InteractionCreate, tag string) {
	response, err := mongodb.Database.FindOne(bson.M{"tag": tag}, mongodb.Database.Collections.Tags)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Tags",
					Description: "Tag not found.",
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})
		return
	}

	content := response["content"].(string)
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}

func createTag(s *discordgo.Session, i *discordgo.InteractionCreate, tag string, content string) {
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
						Title:       "Tags",
						Description: "You do not have permission to create tags.",
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

	_, err := mongodb.Database.InsertOne(bson.M{"tag": tag, "content": content}, mongodb.Database.Collections.Tags)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Tags",
					Description: "Failed to create tag.",
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
				Title:       "Tags",
				Description: fmt.Sprintf("Successfully created tag, `%s`.", tag),
				Color:       0x2B2D31,

				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
				},
			},
		},
	})
}

func deleteTag(s *discordgo.Session, i *discordgo.InteractionCreate, tag string) {
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
						Title:       "Tags",
						Description: "You do not have permission to create tags.",
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

	_, err := mongodb.Database.DeleteOne(bson.M{"tag": tag}, mongodb.Database.Collections.Tags)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Tags",
					Description: "Failed to delete tag.",
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
				Title:       "Tags",
				Description: fmt.Sprintf("Successfully deleted tag, `%s`.", tag),
				Color:       0x2B2D31,

				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
				},
			},
		},
	})
}

func listTags(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cursor, err := mongodb.Database.Find(bson.M{}, mongodb.Database.Collections.Tags)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Tags",
					Description: "Failed to list tags.",
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})
		return
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Tags",
					Description: "Failed to query tags list.",
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})
		return
	}

	var tags string
	for _, tag := range results {
		tags += fmt.Sprintf("`%s`, ", tag["tag"].(string))
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "Tags",
				Description: strings.TrimSuffix(tags, ", "),
				Color:       0x2B2D31,

				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
				},
			},
		},
	})
}
