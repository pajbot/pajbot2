package pkg

import "golang.org/x/oauth2"

type Sender interface {
	TwitchAccount() TwitchAccount
	GetTokenSource() oauth2.TokenSource

	Connected() bool

	Say(Channel, string)
	Mention(Channel, User, string)
	Whisper(User, string)
	Whisperf(User, string, ...interface{})

	// Timeout times the user out a single time immediately
	Timeout(Channel, User, int, string)

	Ban(Channel, User, string)

	GetPoints(Channel, string) uint64

	// give or remove points from user in channel
	BulkEdit(string, []string, int32)

	AddPoints(Channel, string, uint64) (bool, uint64)
	RemovePoints(Channel, string, uint64) (bool, uint64)
	ForceRemovePoints(Channel, string, uint64) uint64

	PointRank(Channel, string) uint64

	// ChannelIDs returns a slice of the channels this bot is connected to
	ChannelIDs() []string

	InChannel(string) bool
	InChannelName(string) bool
	GetUserStore() UserStore
	GetUserContext() UserContext

	GetBotChannel(channelName string) BotChannel
	GetBotChannelByID(channelID string) BotChannel

	MakeUser(string) User
	MakeChannel(string) Channel

	// Permanently join channel with the given channel ID
	JoinChannel(channelID string) error

	// Permanently leave channel with the given channel ID
	LeaveChannel(channelID string) error

	// Connect to the OnNewChannelJoined callback
	OnNewChannelJoined(cb func(channel Channel))

	Quit(message string)

	Application() Application

	// DEV
	Disconnect()
}
