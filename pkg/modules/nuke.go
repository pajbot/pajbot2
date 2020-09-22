package modules

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

func init() {
	Register("nuke", func() pkg.ModuleSpec {
		return &Spec{
			id:    "nuke",
			name:  "Nuke",
			maker: newNuke,

			enabledByDefault: true,
		}
	})
}

const garbageCollectionInterval = 1 * time.Minute
const maxMessageAge = 5 * time.Minute

type nukeModule struct {
	mbase.Base

	messages      map[string][]nukeMessage
	messagesMutex sync.Mutex

	commands pkg.CommandsManager

	ticker *time.Ticker
}

type nukeMessage struct {
	user      pkg.User
	message   pkg.Message
	timestamp time.Time
}

func newNuke(b *mbase.Base) pkg.Module {
	m := &nukeModule{
		Base: *b,

		messages: make(map[string][]nukeMessage),

		commands: commands.NewCommands(),
	}

	m.commands.Register([]string{"!nuke"}, m)

	m.ticker = time.NewTicker(garbageCollectionInterval)

	go func() {
		for range m.ticker.C {
			m.garbageCollect()
		}
	}()

	return m
}

func (m *nukeModule) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !(event.User.IsModerator() || event.User.HasChannelPermission(m.BotChannel().Channel(), pkg.PermissionModeration)) {
		return nil
	}

	if len(parts) < 3 {
		return twitchactions.Mention(event.User, "usage: !nuke bad phrase 1m 10m")
	}
	phrase := strings.Join(parts[1:len(parts)-2], " ")
	scrollbackLength, err := time.ParseDuration(parts[len(parts)-2])
	if err != nil {
		return twitchactions.Mention(event.User, "usage: !nuke bad phrase 1m 10m")
	}
	if scrollbackLength < 0 {
		return twitchactions.Mention(event.User, "usage: !nuke bad phrase 1m 10m")
	}
	timeoutDuration, err := time.ParseDuration(parts[len(parts)-1])
	if err != nil {
		return twitchactions.Mention(event.User, "usage: !nuke bad phrase 1m 10m")
	}
	if timeoutDuration < 0 {
		return twitchactions.Mention(event.User, "usage: !nuke bad phrase 1m 10m")
	}

	m.nuke(event.User, m.BotChannel(), phrase, scrollbackLength, timeoutDuration)

	return nil
}

func (m *nukeModule) OnWhisper(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}

func (m *nukeModule) OnMessage(event pkg.MessageEvent) pkg.Actions {
	defer func() {
		m.addMessage(m.BotChannel().Channel(), event.User, event.Message)
	}()

	return m.commands.OnMessage(event)
}

func (m *nukeModule) garbageCollect() {
	m.messagesMutex.Lock()
	defer m.messagesMutex.Unlock()

	now := time.Now()

	for channelID := range m.messages {
		for i := 0; i < len(m.messages[channelID]); i++ {
			diff := now.Sub(m.messages[channelID][i].timestamp)
			if diff < maxMessageAge {
				m.messages[channelID] = m.messages[channelID][i:]
				break
			}
		}
	}
}

func (m *nukeModule) nuke(source pkg.User, bot pkg.BotChannel, phrase string, scrollbackLength, timeoutDuration time.Duration) {
	if timeoutDuration > 72*time.Hour {
		timeoutDuration = 72 * time.Hour
	}

	lowercasePhrase := strings.ToLower(phrase)

	matcher := func(msg *nukeMessage) bool {
		return strings.Contains(strings.ToLower(msg.message.GetText()), lowercasePhrase)
	}

	reason := "Nuked '" + phrase + "'"

	if strings.HasPrefix(phrase, "/") && strings.HasSuffix(phrase, "/") {
		regex, err := regexp.Compile(phrase[1 : len(phrase)-1])
		if err == nil {
			reason = "Nuked r'" + phrase[1:len(phrase)-1] + "'"
			matcher = func(msg *nukeMessage) bool {
				return regex.MatchString(msg.message.GetText())
			}
		}
		// parse as regex
	}

	now := time.Now()
	timeoutDurationInSeconds := int(timeoutDuration.Seconds())

	if timeoutDurationInSeconds < 1 {
		// Timeout duration too short
		return
	}

	targets := make(map[string]pkg.User)

	m.messagesMutex.Lock()
	defer m.messagesMutex.Unlock()

	messages := m.messages[bot.Channel().GetID()]

	for i := len(messages) - 1; i >= 0; i-- {
		diff := now.Sub(messages[i].timestamp)
		if diff > scrollbackLength {
			// We've gone far enough in the buffer, time to exit
			break
		}

		if matcher(&messages[i]) {
			targets[messages[i].user.GetID()] = messages[i].user
		}
	}

	for _, user := range targets {
		bot.SingleTimeout(user, timeoutDurationInSeconds, reason)
	}

	fmt.Printf("%s nuked %d users for the phrase %s in the last %s for %s\n", source.GetName(), len(targets), phrase, scrollbackLength, timeoutDuration)
}

func (m *nukeModule) addMessage(channel pkg.Channel, user pkg.User, message pkg.Message) {
	m.messagesMutex.Lock()
	defer m.messagesMutex.Unlock()

	m.messages[channel.GetID()] = append(m.messages[channel.GetID()], nukeMessage{
		user:      user,
		message:   message,
		timestamp: time.Now(),
	})
}
