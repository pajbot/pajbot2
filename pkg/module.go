package pkg

import (
	twitch "github.com/gempir/go-twitch-irc"
)

type Module interface {
	Name() string
	Register() error
	OnMessage(channel string, user User, message twitch.Message) error
}
