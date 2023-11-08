package verify

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson"
	"ogusers-bot/internal/config"
	"ogusers-bot/internal/discord"
	"ogusers-bot/internal/ogusers"
	"ogusers-bot/pkg/mongodb"
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

	response, err := ogusers.GetUserByUsername(options[0].Value.(string))
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Verify",
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
					Title:       "Verify",
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

	_, err = mongodb.Database.FindOne(bson.M{"onsite": response.Uid, "verified": true}, mongodb.Database.Collections.Members)
	if err == nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Verify",
					Description: fmt.Sprintf("%s is already linked to another discord account.", response.Username),
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})
		return
	}

	reply, err := mongodb.Database.FindOne(bson.M{
		"discord": i.Member.User.ID,
	}, mongodb.Database.Collections.Members)
	if err != nil {
		mongodb.Database.InsertOne(bson.M{
			"discord":  i.Member.User.ID,
			"onsite":   response.Uid,
			"verified": false,
		}, mongodb.Database.Collections.Members)
	}

	reply, _ = mongodb.Database.FindOne(bson.M{
		"discord": i.Member.User.ID,
	}, mongodb.Database.Collections.Members)

	if reply["verified"].(bool) == true {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Verify",
					Description: "You are already verified.",
					Color:       0x2B2D31,

					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
					},
				},
			},
		})
		return
	}

	verifyCode := fmt.Sprintf("DISCORD-%v", strings.ToUpper(randstr.Hex(25)))

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "Verify",
				Description: fmt.Sprintf("Hey %v, please follow the instructions below to verify your account.", response.Username),
				Color:       0x2B2D31,

				Image: &discordgo.MessageEmbedImage{
					URL: "https://cdn.discordapp.com/attachments/1124116968908275845/1124198721123532911/image.png",
				},

				Footer: &discordgo.MessageEmbedFooter{
					Text: "DO NOT DISMISS THIS MESSAGE",
				},

				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Verification Code",
						Value:  fmt.Sprintf("`%v`", verifyCode),
						Inline: true,
					},
					{
						Name:   "Expires in",
						Value:  fmt.Sprintf("<t:%v:R>", time.Now().Add(time.Minute*5).Unix()),
						Inline: true,
					},
					{
						Name: "Verification Instructions",
						Value: "1. Visit the [signature settings page](https://ogusers.gg/settings.php?settings=signature)\n" +
							"2. Paste your verification code in your signature\n" +
							"3. Wait for the bot to verify your account",
						Inline: false,
					},
				},
			},
		},
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Signature settings",
						URL:   "https://ogusers.gg/settings.php?settings=signature",
						Style: discordgo.LinkButton,
					},
					discordgo.Button{
						Label:    "Copy code",
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("verify-%v", verifyCode),
					},
				},
			},
		},
	})
	discord.C.AddHandler(discord.Component{
		CustomID: fmt.Sprintf("verify-%v", verifyCode),
		Callable: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: verifyCode,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		},
	})

	go func(uid, code string) {
		ticker := time.NewTicker(time.Second * 5)
		done := make(chan bool)

		components := &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Copy code",
						Style:    discordgo.SecondaryButton,
						CustomID: "idfk",
						Disabled: true,
					},
					discordgo.Button{
						Label:    "Signature settings",
						URL:      "https://ogusers.gg/settings.php?settings=signature",
						Style:    discordgo.LinkButton,
						Disabled: true,
					},
				},
			},
		}

		go func() {
			for {
				select {
				case <-done:
					s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
						Components: components,
						Embeds: &[]*discordgo.MessageEmbed{
							{
								Title:       "Verify",
								Description: "Verification expired.",
								Color:       0x2B2D31,

								Thumbnail: &discordgo.MessageEmbedThumbnail{
									URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
								},
							},
						},
					})
					return
				case <-ticker.C:
					response, err := ogusers.GetUserByUID(uid)
					if err != nil {
						s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
							Components: components,
							Embeds: &[]*discordgo.MessageEmbed{
								{
									Title:       "Verify",
									Description: "Failed to get user profile.",
									Color:       0x2B2D31,

									Thumbnail: &discordgo.MessageEmbedThumbnail{
										URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
									},
								},
							},
						})
						continue
					}

					if strings.Contains(response.Signature, code) {
						if _, err := mongodb.Database.FindOneAndUpdate(bson.M{
							"discord": i.Member.User.ID,
						}, bson.M{
							"$set": bson.M{
								"verified": true,
							},
						}, mongodb.Database.Collections.Members); err != nil {
							s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
								Components: components,
								Embeds: &[]*discordgo.MessageEmbed{
									{
										Title:       "Verify",
										Description: "Failed to set database record.",
										Color:       0x2B2D31,

										Thumbnail: &discordgo.MessageEmbedThumbnail{
											URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
										},
									},
								},
							})
							return
						}

						if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, config.C.Roles.Verified); err != nil {
							s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
								Components: components,
								Embeds: &[]*discordgo.MessageEmbed{
									{
										Title:       "Verify",
										Description: "Failed to add verified role, please contact a staff member.",
										Color:       0x2B2D31,

										Thumbnail: &discordgo.MessageEmbedThumbnail{
											URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
										},
									},
								},
							})
							return
						}
						if response.Luminary {
							if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, config.C.Roles.Luminary); err != nil {
								s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
									Components: components,
									Embeds: &[]*discordgo.MessageEmbed{
										{
											Title:       "Verify",
											Description: "Failed to add Luminary role, please contact a staff member.",
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
						if response.Revolution {
							if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, config.C.Roles.Revolution); err != nil {
								s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
									Components: components,
									Embeds: &[]*discordgo.MessageEmbed{
										{
											Title:       "Verify",
											Description: "Failed to add Revolution role, please contact a staff member.",
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

						s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
							Components: components,
							Embeds: &[]*discordgo.MessageEmbed{
								{
									Title:       "Verify",
									Description: "Finished verification process.",
									Color:       0x2B2D31,

									Thumbnail: &discordgo.MessageEmbedThumbnail{
										URL: "https://cdn.discordapp.com/icons/1115373459703353364/34525f9332b1435f86711a4a66818243.png",
									},
								},
							},
						})
						return
					} else {
						continue
					}
				}
			}
		}()

	}(response.Uid, verifyCode)
}
