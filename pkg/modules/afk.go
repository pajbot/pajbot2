package modules

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pajbot/botsync/pkg/client"
	"github.com/pajbot/botsync/pkg/protocol"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/utils"
)

func init() {
	Register("afk", func() pkg.ModuleSpec {
		afkDatabase := map[string]bool{}

		return &moduleSpec{
			id:               "afk",
			name:             "AFK",
			enabledByDefault: false,
			maker: func(b base) pkg.Module {
				return newAFK(b, afkDatabase)
			},
		}
	})
}

type afk struct {
	base

	commands map[string]pkg.CustomCommand

	botsync *client.Client

	afkDatabase map[string]bool
}

func newAFK(b base, afkDatabase map[string]bool) pkg.Module {
	m := &afk{
		base: b,

		afkDatabase: afkDatabase,

		commands: make(map[string]pkg.CustomCommand),

		botsync: client.NewClient("ws://localhost:8080/ws/pubsub"),
	}

	m.botsync.OnMessage(m.onMessage)

	m.botsync.OnConnect(func() {
		fmt.Println("Connected to botsync")
		m.botsync.Send(protocol.NewAFKSubscribeMessage(m.bot.ChannelID()))
		m.botsync.Send(protocol.NewBackSubscribeMessage(m.bot.ChannelID()))
	})

	// FIXME
	m.Initialize()

	return m
}

type afkCmd struct {
	m *afk
}

func (c *afkCmd) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if user.IsModerator() || user.IsBroadcaster(channel) {
		c.m.botsync.Send(protocol.NewAFKMessage(&protocol.AFKParameters{
			UserID:   user.GetID(),
			UserName: user.GetName(),

			Reason: strings.Join(parts[1:], " "),

			ChannelID:   botChannel.ChannelID(),
			ChannelName: botChannel.ChannelName(),
		}))
	}
}

func (m *afk) registerCommand(aliases []string, command pkg.CustomCommand) {
	for _, alias := range aliases {
		m.commands[alias] = command
	}
}

func (m *afk) Initialize() {
	m.registerCommand([]string{"!afk", "!gn"}, &afkCmd{m})

	m.botsync.SetAuthentication(protocol.Authentication{
		TwitchUserID:        m.bot.Bot().TwitchAccount().ID(),
		AuthenticationToken: "penis",
	})

	go func() {
		for {
			err := m.botsync.Connect()
			if err != nil {
				if err == client.ErrDisconnected {
					return
				}
			}

			// wait a second before attempting to reconnect
			<-time.After(time.Second)
		}
	}()
}

func (m *afk) onMessage(message *protocol.Message) {
	switch message.Topic {
	case "afk":
		parameters := protocol.AFKParameters{}
		err := json.Unmarshal(message.Data, &parameters)
		if err != nil {
			fmt.Println("Error parsing afk message:", err)
			return
		}
		if !message.Historic {
			m.bot.Say(parameters.UserName + " just went afk: " + parameters.Reason)
		}
		m.afkDatabase[parameters.UserID] = true

	case "back":
		parameters := protocol.BackParameters{}
		err := json.Unmarshal(message.Data, &parameters)
		if err != nil {
			fmt.Println("Error parsing afk message:", err)
			return
		}
		afkDuration := time.Millisecond * time.Duration(parameters.Duration)
		response := fmt.Sprintf("%s just came back after %s: %s",
			parameters.UserName, utils.DurationString(afkDuration), parameters.Reason)
		m.bot.Say(response)
		delete(m.afkDatabase, parameters.UserID)
	}
}

func (m *afk) Disable() error {
	m.botsync.Disconnect()

	return m.base.Disable()
}

func (m *afk) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	if _, ok := m.afkDatabase[user.GetID()]; ok {
		m.botsync.Send(protocol.NewBackMessage(&protocol.BackParameters{
			UserID:   user.GetID(),
			UserName: user.GetName(),

			ChannelID:   bot.ChannelID(),
			ChannelName: bot.ChannelName(),
		}))
		return nil
	}

	parts := strings.Split(message.GetText(), " ")
	if len(parts) == 0 {
		return nil
	}

	if command, ok := m.commands[strings.ToLower(parts[0])]; ok {
		command.Trigger(m.bot, parts, bot.Channel(), user, message, action)
	}

	return nil
}
