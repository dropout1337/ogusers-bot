package forceverify

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slices"
	"ogusers-bot/internal/config"
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
						Title:       "Force verify",
						Description: "You do not have permission to force verify.",
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

	reply, err := mongodb.Database.FindOne(bson.M{"onsite": options[1].Value.(string), "verified": true}, mongodb.Database.Collections.Members)
	if err == nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Force verify",
					Description: fmt.Sprintf("UID **`%s`** is linked to a <@%s> (`%s`).", options[1].Value.(string), reply["discord"].(string), reply["discord"].(string)),
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})
		return
	}

	mongodb.Database.InsertOne(bson.M{"onsite": options[1].Value.(string), "discord": options[0].Value.(string), "verified": true}, mongodb.Database.Collections.Members)
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "Force verify",
				Description: fmt.Sprintf("Successfully force verified UID **`%s`** to <@%s> (`%s`).", options[1].Value.(string), options[0].Value.(string), options[0].Value.(string)),
				Color:       0x2B2D31,

				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
				},
			},
		},
	})
}
