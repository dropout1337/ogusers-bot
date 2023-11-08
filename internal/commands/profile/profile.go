package profile

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"ogusers-bot/internal/ogusers"
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
	case "uid":
		response, err := ogusers.GetUserByUID(options[0].Options[0].Value.(string))
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Profile",
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
						Title:       "Profile",
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

		profile(s, i, response)
	case "username":
		response, err := ogusers.GetUserByUsername(options[0].Options[0].Value.(string))
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Profile",
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
						Title:       "Profile",
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

		profile(s, i, response)
	case "discord":
		reply, err := mongodb.Database.FindOne(bson.M{"discord": options[0].Options[0].Value.(string), "verified": true}, mongodb.Database.Collections.Members)
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Profile",
						Description: "User has no OGUsers account linked.",
						Color:       0x2B2D31,

						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
						},
					},
				},
			})
			return
		}

		response, err := ogusers.GetUserByUID(reply["onsite"].(string))
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Profile",
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
						Title:       "Profile",
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

		profile(s, i, response)
	}
}

func profile(s *discordgo.Session, i *discordgo.InteractionCreate, response *ogusers.Response) {
	var rank string
	if response.Rank == "0" {
		rank = "No rank"
	} else {
		rank = response.Rank
	}

	var discord string
	reply, err := mongodb.Database.FindOne(bson.M{"onsite": response.Uid}, mongodb.Database.Collections.Members)
	if err != nil {
		discord = "Not linked"
	} else {
		discord = fmt.Sprintf("<@%v>", reply["discord"].(string))
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       fmt.Sprintf("%v's profile", strings.ToLower(response.Username)),
				Description: fmt.Sprintf("%v", strings.ReplaceAll(response.Title, "&amp;", "&")),
				URL:         fmt.Sprintf("https://ogusers.gg/account.php?uid=%v", response.Uid),
				Color:       0x2B2D31,

				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: response.Avatar,
				},

				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("UID: %v | Joined: %v | Rank: %v", response.Uid, response.Registered, rank),
				},

				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Reputation",
						Value:  fmt.Sprintf("[%v](https://ogusers.gg/reputation.php?uid=%v)", response.Rep, response.Uid),
						Inline: true,
					},
					{
						Name:   "Vouches",
						Value:  fmt.Sprintf("[%v](https://ogusers.gg/vouches.php?uid=%v)", response.Vouches, response.Uid),
						Inline: true,
					},
					{
						Name:   "Discord",
						Value:  discord,
						Inline: true,
					},
					{
						Name:   "Credits",
						Value:  response.Credits,
						Inline: true,
					},
					{
						Name:   "Items",
						Value:  response.Items,
						Inline: true,
					},
					{
						Name:   "Last online",
						Value:  fmt.Sprintf("<t:%v:R>", response.Lastonline),
						Inline: true,
					},
				},
			},
		},
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Open profile",
						URL:   "https://ogusers.gg/account.php?uid=" + response.Uid,
						Style: discordgo.LinkButton,
					},
					discordgo.Button{
						Label: "Rep",
						URL:   "https://ogusers.gg/reputation.php?uid=" + response.Uid,
						Style: discordgo.LinkButton,
					},
					discordgo.Button{
						Label: "Vouch",
						URL:   "https://ogusers.gg/vouches.php?uid=" + response.Uid,
						Style: discordgo.LinkButton,
					},
				},
			},
		},
	})
}
