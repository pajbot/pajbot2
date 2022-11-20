package nuke

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
	"github.com/pajbot/pajbot2/pkg/modules"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

var (
	ErrUsage = errors.New("usage: !nuke bad phrase 1m 10m")
)

func init() {
	modules.Register("nuke", func() pkg.ModuleSpec {
		return modules.NewSpec("nuke", "Nuke", true, func(b *mbase.Base) pkg.Module {
			return NewNuke(b, &NukeParameterParser{})
		})
	})
}

const garbageCollectionInterval = 1 * time.Minute
const maxMessageAge = 5 * time.Minute

type NukeParameters struct {
	Phrase      string
	RegexPhrase *regexp.Regexp

	ScrollbackLength time.Duration

	TimeoutDuration time.Duration
}

type NukeModule struct {
	mbase.Base

	messages      map[string][]nukeMessage
	messagesMutex sync.Mutex

	parameterParser *NukeParameterParser

	commands pkg.CommandsManager

	ticker *time.Ticker
}

type nukeMessage struct {
	user      pkg.User
	message   pkg.Message
	timestamp time.Time
}

func NewNuke(b *mbase.Base, parameterParser *NukeParameterParser) *NukeModule {
	m := &NukeModule{
		Base: *b,

		messages: make(map[string][]nukeMessage),

		commands: commands.NewCommands(),

		parameterParser: parameterParser,
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

func (m *NukeModule) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !(event.User.IsModerator() || event.User.HasPermission(m.BotChannel().Channel(), pkg.PermissionModeration)) {
		return nil
	}

	params, err := m.parameterParser.ParseNukeParameters(parts)
	if err != nil {
		return twitchactions.Mention(event.User, err.Error())
	}

	m.nuke(event.User, m.BotChannel(), params)

	return nil
}

func (m *NukeModule) OnWhisper(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}

func (m *NukeModule) OnMessage(event pkg.MessageEvent) pkg.Actions {
	defer func() {
		m.addMessage(m.BotChannel().Channel(), event.User, event.Message)
	}()

	return m.commands.OnMessage(event)
}

func (m *NukeModule) garbageCollect() {
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

func (m *NukeModule) nuke(source pkg.User, bot pkg.BotChannel, params *NukeParameters) {
	// TODO: Should this be moved to the NukeParameters parser?
	if params.TimeoutDuration > 72*time.Hour {
		params.TimeoutDuration = 72 * time.Hour
	}

	lowercasePhrase := strings.ToLower(params.Phrase)

	matcher := func(msg *nukeMessage) bool {
		return strings.Contains(strings.ToLower(msg.message.GetText()), lowercasePhrase)
	}

	if params.RegexPhrase != nil {
		matcher = func(msg *nukeMessage) bool {
			return params.RegexPhrase.MatchString(msg.message.GetText())
		}
	}

	reason := "Nuked '" + params.Phrase + "'"

	now := time.Now()
	timeoutDurationInSeconds := int(params.TimeoutDuration.Seconds())

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
		if diff > params.ScrollbackLength {
			// We've gone far enough in the buffer, time to exit
			break
		}

		if matcher(&messages[i]) {
			targets[messages[i].user.GetID()] = messages[i].user
		}
	}

	for _, user := range targets {
		bot.Timeout(user, timeoutDurationInSeconds, reason)
	}

	fmt.Printf("%s nuked %d users for the phrase %s in the last %s for %s\n", source.GetName(), len(targets), params.Phrase, params.ScrollbackLength, params.TimeoutDuration)
}

func (m *NukeModule) addMessage(channel pkg.Channel, user pkg.User, message pkg.Message) {
	m.messagesMutex.Lock()
	defer m.messagesMutex.Unlock()

	m.messages[channel.GetID()] = append(m.messages[channel.GetID()], nukeMessage{
		user:      user,
		message:   message,
		timestamp: time.Now(),
	})
}
