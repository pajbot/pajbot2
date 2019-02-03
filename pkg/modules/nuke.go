package modules

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/pajlada/pajbot2/pkg"
)

const garbageCollectionInterval = 1 * time.Minute
const maxMessageAge = 5 * time.Minute

type nukeModule struct {
	botChannel pkg.BotChannel

	server        *server
	messages      map[string][]nukeMessage
	messagesMutex sync.Mutex

	ticker *time.Ticker
}

type nukeMessage struct {
	user      pkg.User
	message   pkg.Message
	timestamp time.Time
}

func newNuke() pkg.Module {
	m := &nukeModule{
		server:   &_server,
		messages: make(map[string][]nukeMessage),
	}

	m.ticker = time.NewTicker(garbageCollectionInterval)

	go func() {
		for {
			select {
			case <-m.ticker.C:
				m.garbageCollect()
			}
		}
	}()

	return m
}

var nukeSpec = moduleSpec{
	id:    "nuke",
	name:  "Nuke",
	maker: newNuke,

	enabledByDefault: true,
}

func (m *nukeModule) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

	return nil
}

func (m *nukeModule) Disable() error {
	return nil
}

func (m *nukeModule) Spec() pkg.ModuleSpec {
	return &nukeSpec
}

func (m *nukeModule) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m *nukeModule) OnWhisper(bot pkg.BotChannel, user pkg.User, message pkg.Message) error {
	const usageString = `Usage: #channel !nuke phrase phrase phrase time`

	parts := strings.Split(message.GetText(), " ")
	// Minimum required parts: 4
	// !nuke PHRASE SCROLLBACK_LENGTH TIMEOUT_DURATION
	if len(parts) >= 4 {
		if parts[0] != "!nuke" {
			return nil
		}

		// TODO: Add another specific global/channel permission to check
		if !user.IsModerator() && !user.IsBroadcaster(bot.Channel()) && !user.HasChannelPermission(bot.Channel(), pkg.PermissionModeration) && !user.HasGlobalPermission(pkg.PermissionModeration) {
			return nil
		}

		phrase := strings.Join(parts[1:len(parts)-2], " ")
		scrollbackLength, err := time.ParseDuration(parts[len(parts)-2])
		if err != nil {
			bot.Mention(user, "usage: !nuke bad phrase 1m 10m")
			return err
		}
		if scrollbackLength < 0 {
			bot.Mention(user, "usage: !nuke bad phrase 1m 10m")
			return errors.New("scrollback length must be positive")
		}
		timeoutDuration, err := time.ParseDuration(parts[len(parts)-1])
		if err != nil {
			bot.Mention(user, "usage: !nuke bad phrase 1m 10m")
			return err
		}
		if timeoutDuration < 0 {
			bot.Mention(user, "usage: !nuke bad phrase 1m 10m")
			return errors.New("timeout duration must be positive")
		}

		m.nuke(user, bot, phrase, scrollbackLength, timeoutDuration)
	}

	return nil
}

func (m *nukeModule) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	defer func() {
		m.addMessage(bot.Channel(), user, message)
	}()

	parts := strings.Split(message.GetText(), " ")
	// Minimum required parts: 4
	// !nuke PHRASE SCROLLBACK_LENGTH TIMEOUT_DURATION
	if len(parts) >= 4 {
		if parts[0] != "!nuke" {
			return nil
		}

		// TODO: Add another specific global/channel permission to check
		if !user.IsModerator() && !user.IsBroadcaster(bot.Channel()) && !user.HasChannelPermission(bot.Channel(), pkg.PermissionModeration) && !user.HasGlobalPermission(pkg.PermissionModeration) {
			return nil
		}

		phrase := strings.Join(parts[1:len(parts)-2], " ")
		scrollbackLength, err := time.ParseDuration(parts[len(parts)-2])
		if err != nil {
			bot.Mention(user, "usage: !nuke bad phrase 1m 10m")
			return err
		}
		if scrollbackLength < 0 {
			bot.Mention(user, "usage: !nuke bad phrase 1m 10m")
			return errors.New("scrollback length must be positive")
		}
		timeoutDuration, err := time.ParseDuration(parts[len(parts)-1])
		if err != nil {
			bot.Mention(user, "usage: !nuke bad phrase 1m 10m")
			return err
		}
		if timeoutDuration < 0 {
			bot.Mention(user, "usage: !nuke bad phrase 1m 10m")
			return errors.New("timeout duration must be positive")
		}

		m.nuke(user, bot, phrase, scrollbackLength, timeoutDuration)
	}

	return nil
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
	if timeoutDuration > 24*time.Hour {
		timeoutDuration = 24 * time.Hour
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
		bot.Timeout(user, timeoutDurationInSeconds, reason)
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
