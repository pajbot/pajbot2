package twitchactions

import "github.com/pajbot/pajbot2/pkg"

type Whisper struct {
	user    pkg.User
	content string
}

func (w Whisper) User() pkg.User {
	return w.user
}

func (w Whisper) Content() string {
	return w.content
}
