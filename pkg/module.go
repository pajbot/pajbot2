package pkg

import (
	twitch "github.com/gempir/go-twitch-irc"
)

type Module interface {
	Name() string
	Register() error
	OnWhisper(source User, message twitch.Message) error
	OnMessage(source Channel, user User, message twitch.Message) error
}
