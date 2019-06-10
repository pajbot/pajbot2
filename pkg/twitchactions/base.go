package twitchactions

import (
	"fmt"
	"sync"
	"time"

	"github.com/pajbot/pajbot2/pkg"
)

type baseActions struct {
	mutesMutex sync.RWMutex
	mutes      []pkg.MuteAction

	messagesMutex sync.RWMutex
	messages      []pkg.MessageAction

	whispersMutex sync.RWMutex
	whispers      []pkg.WhisperAction
}

func (a *baseActions) Timeout(user pkg.User, duration time.Duration) pkg.MuteAction {
	action := &Timeout{
		mute: mute{
			user: user,
		},
		duration: duration,
	}

	a.mutesMutex.Lock()
	a.mutes = append(a.mutes, action)
	a.mutesMutex.Unlock()

	return action
}

func (a *baseActions) Ban(user pkg.User) pkg.MuteAction {
	action := &Ban{
		mute: mute{
			user: user,
		},
	}

	a.mutesMutex.Lock()
	a.mutes = append(a.mutes, action)
	a.mutesMutex.Unlock()

	return action
}

func (a *baseActions) Say(content string) pkg.MessageAction {
	action := &Message{
		content: content,
	}

	a.messagesMutex.Lock()
	a.messages = append(a.messages, action)
	a.messagesMutex.Unlock()

	return action
}

func (a *baseActions) Mention(user pkg.User, content string) pkg.MessageAction {
	action := &Message{
		content: fmt.Sprintf("@%s, %s", user.GetName(), content),
	}

	a.messagesMutex.Lock()
	a.messages = append(a.messages, action)
	a.messagesMutex.Unlock()

	return action
}

func (a *baseActions) Whisper(user pkg.User, content string) pkg.WhisperAction {
	action := &Whisper{
		user:    user,
		content: content,
	}

	a.whispersMutex.Lock()
	a.whispers = append(a.whispers, action)
	a.whispersMutex.Unlock()

	return action
}

func (a *baseActions) Mutes() []pkg.MuteAction {
	a.mutesMutex.RLock()
	defer a.mutesMutex.RUnlock()

	return a.mutes
}

func (a *baseActions) Messages() []pkg.MessageAction {
	a.messagesMutex.RLock()
	defer a.messagesMutex.RUnlock()

	return a.messages
}

func (a *baseActions) Whispers() []pkg.WhisperAction {
	a.whispersMutex.RLock()
	defer a.whispersMutex.RUnlock()

	return a.whispers
}
