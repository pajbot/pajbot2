package pkg

import "golang.org/x/oauth2"

type TwitchAuths interface {
	Bot() *oauth2.Config
	Streamer() *oauth2.Config
	User() *oauth2.Config
}
