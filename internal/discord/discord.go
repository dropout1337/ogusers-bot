package discord

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
)

var C *Discord

func New(token string) (*Discord, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	client := Discord{
		Session: session,
	}

	client.Commands = make(map[string]func(*discordgo.Session, *discordgo.InteractionCreate), 0)
	client.MessageComponents = make(map[string]func(*discordgo.Session, *discordgo.InteractionCreate), 0)

	return &client, nil
}

func (d *Discord) AddHandler(v any) error {
	switch v.(type) {
	case Command:
		cmd := v.(Command)
		d.Commands[cmd.ApplicationCommand.Name] = cmd.Callable

		command, err := d.Session.ApplicationCommandCreate(d.Session.State.User.ID, "", cmd.ApplicationCommand)
		if err != nil {
			return err
		} else {
			d.RegisteredCommands = append(d.RegisteredCommands, command)
		}
	case Component:
		component := v.(Component)
		d.MessageComponents[component.CustomID] = component.Callable
	default:
		return errors.New("invalid handler type")
	}

	return nil
}

func (d *Discord) Open() error {
	err := d.Session.Open()
	if err != nil {
		return err
	}

	d.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:

			if callable, ok := d.Commands[i.ApplicationCommandData().Name]; ok {
				callable(s, i)
			}

		case discordgo.InteractionMessageComponent:
			if callable, ok := d.MessageComponents[i.MessageComponentData().CustomID]; ok {
				callable(s, i)
			}
		case discordgo.InteractionModalSubmit:
			if callable, ok := d.MessageComponents[i.ModalSubmitData().CustomID]; ok {
				callable(s, i)
			}
		}
	})

	defer d.Session.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	for _, v := range d.RegisteredCommands {
		err := d.Session.ApplicationCommandDelete(d.Session.State.User.ID, "", v.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
