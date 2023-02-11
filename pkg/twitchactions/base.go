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

	unmutesMutex sync.RWMutex
	unmutes      []pkg.UnmuteAction

	deletesMutex sync.RWMutex
	deletes      []pkg.DeleteAction

	messagesMutex sync.RWMutex
	messages      []pkg.MessageAction

	whispersMutex sync.RWMutex
	whispers      []pkg.WhisperAction
}

func (a *baseActions) StopPropagation() bool {
	a.mutesMutex.Lock()
	defer a.mutesMutex.Unlock()
	return len(a.mutes) > 0
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

func (a *baseActions) Unban(user pkg.User) pkg.UnmuteAction {
	action := &unmute{
		user:     user,
		muteType: pkg.MuteTypePermanent,
	}

	a.unmutesMutex.Lock()
	a.unmutes = append(a.unmutes, action)
	a.unmutesMutex.Unlock()

	return action
}

func (a *baseActions) Untimeout(user pkg.User) pkg.UnmuteAction {
	action := &unmute{
		user:     user,
		muteType: pkg.MuteTypeTemporary,
	}

	a.unmutesMutex.Lock()
	a.unmutes = append(a.unmutes, action)
	a.unmutesMutex.Unlock()

	return action
}

func (a *baseActions) Delete(message string) pkg.DeleteAction {
	action := &deleteAction{
		message: message,
	}

	a.deletesMutex.Lock()
	a.deletes = append(a.deletes, action)
	a.deletesMutex.Unlock()

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

func (a *baseActions) Unmutes() []pkg.UnmuteAction {
	a.unmutesMutex.RLock()
	defer a.unmutesMutex.RUnlock()

	return a.unmutes
}

func (a *baseActions) Deletes() []pkg.DeleteAction {
	a.deletesMutex.RLock()
	defer a.deletesMutex.RUnlock()

	return a.deletes
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
