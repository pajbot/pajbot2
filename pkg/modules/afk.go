package modules

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pajbot/botsync/pkg/client"
	"github.com/pajbot/botsync/pkg/protocol"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/utils"
)

func init() {
	Register("afk", func() pkg.ModuleSpec {
		afkDatabase := map[string]bool{}

		return &moduleSpec{
			id:               "afk",
			name:             "AFK",
			enabledByDefault: false,
			maker: func(b mbase.Base) pkg.Module {
				return newAFK(b, afkDatabase)
			},
		}
	})
}

type afk struct {
	mbase.Base

	commands pkg.CommandsManager

	botsync *client.Client

	afkDatabase map[string]bool
}

func newAFK(b mbase.Base, afkDatabase map[string]bool) pkg.Module {
	m := &afk{
		Base: b,

		afkDatabase: afkDatabase,

		commands: commands.NewCommands(),

		botsync: client.NewClient("ws://localhost:8080/ws/pubsub"),
	}

	m.botsync.OnMessage(m.onMessage)

	m.botsync.OnConnect(func() {
		fmt.Println("Connected to botsync")
		m.botsync.Send(protocol.NewAFKSubscribeMessage(m.BotChannel().ChannelID()))
		m.botsync.Send(protocol.NewBackSubscribeMessage(m.BotChannel().ChannelID()))
	})

	// FIXME
	m.Initialize()

	return m
}

type afkCmd struct {
	m *afk
}

func (c *afkCmd) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	// Temporarily limit afk commands to moderators, since it can be used to say stupid things
	if !event.User.IsModerator() {
		return nil
	}

	c.m.botsync.Send(protocol.NewAFKMessage(&protocol.AFKParameters{
		UserID:   event.User.GetID(),
		UserName: event.User.GetName(),

		Reason: strings.Join(parts[1:], " "),

		ChannelID:   event.Channel.GetID(),
		ChannelName: event.Channel.GetName(),
	}))

	return nil
}

func (m *afk) Initialize() {
	m.commands.Register([]string{"!afk", "!gn"}, &afkCmd{m})

	m.botsync.SetAuthentication(protocol.Authentication{
		TwitchUserID:        m.BotChannel().Bot().TwitchAccount().ID(),
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
			m.BotChannel().Say(parameters.UserName + " just went afk: " + parameters.Reason)
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
		m.BotChannel().Say(response)
		delete(m.afkDatabase, parameters.UserID)
	}
}

func (m *afk) Disable() error {
	m.botsync.Disconnect()

	return m.Base.Disable()
}

func (m *afk) OnMessage(event pkg.MessageEvent) pkg.Actions {
	user := event.User

	if _, ok := m.afkDatabase[user.GetID()]; ok {
		m.botsync.Send(protocol.NewBackMessage(&protocol.BackParameters{
			UserID:   user.GetID(),
			UserName: user.GetName(),

			ChannelID:   m.BotChannel().ChannelID(),
			ChannelName: m.BotChannel().ChannelName(),
		}))
		return nil
	}

	return m.commands.OnMessage(event)
}
