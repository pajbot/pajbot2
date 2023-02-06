package pkg

import (
	"time"
)

type MuteType uint

const (
	MuteTypeTemporary MuteType = iota
	MuteTypePermanent
)

// MuteAction defines an action that will mute/timeout/ban or otherwise stop a user from participating in chat, either temporarily or permanently
type MuteAction interface {
	User() User
	SetReason(reason string)
	Reason() string

	Type() MuteType

	Duration() time.Duration
}

// UnmuteAction defines an action that will unmute/untimeout or unban a user
type UnmuteAction interface {
	User() User

	Type() MuteType
}

// DeleteAction defines an action that will delete a message
type DeleteAction interface {
	Message() string
}

// MessageAction defines a message that will be publicly displayed
type MessageAction interface {
	// TODO: Add reply message action
	SetAction(v bool)
	Evaluate() string
}

// WhisperAction defines a message that will be privately sent to a user
type WhisperAction interface {
	User() User
	Content() string
}

// Actions is a list of actions that wants to be run
// An implementation of this can decide to filter out all mutes except for the "most grave one"
type Actions interface {
	Timeout(user User, duration time.Duration) MuteAction

	Ban(user User) MuteAction

	Unban(user User) UnmuteAction

	Say(content string) MessageAction

	Delete(message string) DeleteAction

	Mention(user User, content string) MessageAction

	Whisper(user User, content string) WhisperAction

	Mutes() []MuteAction
	Unmutes() []UnmuteAction
	Deletes() []DeleteAction
	Messages() []MessageAction
	Whispers() []WhisperAction

	StopPropagation() bool

	// DoOnSuccess(func())

	// Do(func())
}
