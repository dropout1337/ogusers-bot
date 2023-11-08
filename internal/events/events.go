package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"ogusers-bot/internal/config"
	"ogusers-bot/pkg/logging"
	"ogusers-bot/pkg/mongodb"
	"strings"
	"time"
)

func ChannelPurger(s *discordgo.Session) {
	ticker := time.NewTicker(6 * time.Hour)

	go func() {
		for {
			select {
			case <-ticker.C:
				channels, err := s.GuildChannels(config.C.Guild.Id)
				if err != nil {
					logging.Logger.Error().
						Err(err).
						Msg("failed to get channels")
					continue
				}

				var generalChannel *discordgo.Channel
				for _, channel := range channels {
					if channel.Name == config.C.Guild.GeneralChannel {
						generalChannel = channel
					}
				}
				if generalChannel == nil {
					logging.Logger.Error().
						Err(err).
						Msg("failed to get channel")
					continue
				}

				s.ChannelMessageSendEmbed(generalChannel.ID, &discordgo.MessageEmbed{
					Title:       "Channel Purger",
					Description: fmt.Sprintf("This channel will be deleted in <t:%v:R>.", time.Now().Add(15*time.Second).Unix()),
					Color:       0x2B2D31,
				})

				time.Sleep(15 * time.Second)

				_, err = s.ChannelDelete(generalChannel.ID)
				if err != nil {
					logging.Logger.Error().
						Err(err).
						Msg("failed to delete channel")
					continue
				}

				response, err := s.GuildChannelCreateComplex(config.C.Guild.Id, discordgo.GuildChannelCreateData{
					Name:                 generalChannel.Name,
					Type:                 generalChannel.Type,
					Topic:                generalChannel.Topic,
					NSFW:                 generalChannel.NSFW,
					Position:             generalChannel.Position,
					Bitrate:              generalChannel.Bitrate,
					UserLimit:            generalChannel.UserLimit,
					RateLimitPerUser:     generalChannel.RateLimitPerUser,
					ParentID:             generalChannel.ParentID,
					PermissionOverwrites: generalChannel.PermissionOverwrites,
				})
				if err != nil {
					logging.Logger.Error().
						Err(err).
						Msg("failed to create channel")
					continue
				}

				s.GuildChannelsReorder(config.C.Guild.Id, []*discordgo.Channel{
					{
						ID:       response.ID,
						Position: generalChannel.Position,
					},
				})

				s.ChannelMessageSendEmbed(response.ID, &discordgo.MessageEmbed{
					Title:       "Channel Purger",
					Description: fmt.Sprintf("Successfully purged <#%v> at <t:%v>.", response.ID, time.Now().Unix()),
					Color:       0x2B2D31,

					Image: &discordgo.MessageEmbedImage{
						URL: "https://media.tenor.com/M5xHAKQRpNYAAAAC/bh187-spongebob.gif",
					},
				})
				msg, err := s.ChannelMessageSend(response.ID, fmt.Sprintf("<@&%v>", config.C.Roles.PurgeRevive))
				if err != nil {
					logging.Logger.Error().
						Err(err).
						Msg("failed to send message")
					continue
				} else {
					time.Sleep(500 * time.Millisecond)
					s.ChannelMessageDelete(response.ID, msg.ID)
				}

				logging.Logger.Info().
					Str("channel", response.ID).
					Msg("channel nuked")
			}
		}
	}()
}

func OnMessage(s *discordgo.Session, r *discordgo.MessageCreate) {
	if r.GuildID != config.C.Guild.Id {
		return
	}

	go func() {
		if channel, err := s.Channel(r.ChannelID); err != nil {
			return
		} else if channel.Name != config.C.Guild.GeneralChannel {
			return
		}

		var files []*discordgo.File
		for _, attachment := range r.Attachments {
			resp, err := http.Get(attachment.URL)
			if err != nil {
				logging.Logger.Error().
					Err(err).
					Msg("failed to get attachment")
				continue
			}

			files = append(files, &discordgo.File{
				Name:        attachment.Filename,
				ContentType: attachment.ContentType,
				Reader:      resp.Body,
			})
		}

		var safeContent string
		safeContent = strings.ReplaceAll(r.Content, "@everyone", "@\u200beveryone")
		safeContent = strings.ReplaceAll(r.Content, "@", "@\u200b")

		s.WebhookExecute(strings.Split(config.C.Guild.MirrorWebhook, "/")[0], strings.Split(config.C.Guild.MirrorWebhook, "/")[1], false, &discordgo.WebhookParams{
			Content:   safeContent,
			Username:  r.Author.Username,
			AvatarURL: r.Author.AvatarURL(""),
			Files:     files,
			Embeds:    r.Embeds,
		})
	}()

	if strings.HasPrefix(r.Content, "og.") {
		tagName := strings.Split(r.Content, "og.")[1]
		reply, err := mongodb.Database.FindOne(bson.M{"tag": tagName}, mongodb.Database.Collections.Tags)
		if err != nil {
			return
		}

		s.ChannelMessageSend(r.ChannelID, reply["content"].(string))
	}
}
