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
	Register(afkSpec)
}

var (
	afkDatabase = map[string]bool{}
)

type afk struct {
	botChannel pkg.BotChannel

	server *server

	commands map[string]pkg.CustomCommand

	botsync *client.Client
}

var afkSpec = &moduleSpec{
	id:               "afk",
	name:             "AFK",
	maker:            newAFK,
	enabledByDefault: false,
}

func newAFK() pkg.Module {
	m := &afk{
		server: &_server,

		commands: make(map[string]pkg.CustomCommand),

		botsync: client.NewClient("ws://localhost:8080/ws/pubsub"),
	}

	m.botsync.OnMessage(m.onMessage)

	m.botsync.OnConnect(func() {
		fmt.Println("Connected to botsync")
		m.botsync.Send(protocol.NewAFKSubscribeMessage(m.botChannel.ChannelID()))
		m.botsync.Send(protocol.NewBackSubscribeMessage(m.botChannel.ChannelID()))
	})

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

func (m *afk) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

	m.registerCommand([]string{"!afk", "!gn"}, &afkCmd{m})

	m.botsync.SetAuthentication(protocol.Authentication{
		TwitchUserID:        botChannel.Bot().TwitchAccount().ID(),
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

	return nil
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
			m.botChannel.Say(parameters.UserName + " just went afk: " + parameters.Reason)
		}
		afkDatabase[parameters.UserID] = true

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
		m.botChannel.Say(response)
		delete(afkDatabase, parameters.UserID)
	}
}

func (m *afk) Disable() error {
	m.botsync.Disconnect()
	return nil
}

func (m *afk) Spec() pkg.ModuleSpec {
	return afkSpec
}

func (m *afk) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m *afk) OnWhisper(bot pkg.BotChannel, source pkg.User, message pkg.Message) error {
	return nil
}

func (m *afk) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	if _, ok := afkDatabase[user.GetID()]; ok {
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
		command.Trigger(m.botChannel, parts, bot.Channel(), user, message, action)
	}

	return nil
}
