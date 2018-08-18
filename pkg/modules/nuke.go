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

var _ pkg.Module = &Nuke{}

const garbageCollectionInterval = 1 * time.Minute
const maxMessageAge = 5 * time.Minute

type nukeMessage struct {
	channel   pkg.Channel
	user      pkg.User
	message   pkg.Message
	timestamp time.Time
}

type Nuke struct {
	server        *server
	messages      []nukeMessage
	messagesMutex sync.Mutex

	ticker *time.Ticker
}

func NewNuke() *Nuke {
	m := &Nuke{
		server: &_server,
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

func (m *Nuke) Register() error {
	return nil
}

func (m *Nuke) Name() string {
	return "Nuke"
}

func (m *Nuke) OnWhisper(bot pkg.Sender, user pkg.User, message pkg.Message) error {
	return nil
}

func (m *Nuke) OnMessage(bot pkg.Sender, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	defer func() {
		m.addMessage(channel, user, message)
	}()

	parts := strings.Split(message.GetText(), " ")
	// Minimum required parts: 4
	// !nuke PHRASE SCROLLBACK_LENGTH TIMEOUT_DURATION
	if len(parts) >= 4 {
		if parts[0] != "!nuke" {
			return nil
		}

		// TODO: Add another specific global/channel permission to check
		if !user.IsModerator() && !user.IsBroadcaster(channel) && !user.HasChannelPermission(channel, pkg.PermissionModeration) && !user.HasGlobalPermission(pkg.PermissionModeration) {
			return nil
		}

		phrase := strings.Join(parts[1:len(parts)-2], " ")
		scrollbackLength, err := time.ParseDuration(parts[len(parts)-2])
		if err != nil {
			bot.Mention(channel, user, "usage: !nuke bad phrase 1m 10m")
			return err
		}
		if scrollbackLength < 0 {
			bot.Mention(channel, user, "usage: !nuke bad phrase 1m 10m")
			return errors.New("scrollback length must be positive")
		}
		timeoutDuration, err := time.ParseDuration(parts[len(parts)-1])
		if err != nil {
			bot.Mention(channel, user, "usage: !nuke bad phrase 1m 10m")
			return err
		}
		if timeoutDuration < 0 {
			bot.Mention(channel, user, "usage: !nuke bad phrase 1m 10m")
			return errors.New("timeout duration must be positive")
		}

		m.nuke(user, bot, channel, phrase, scrollbackLength, timeoutDuration)
	}

	return nil
}

func (m *Nuke) garbageCollect() {
	m.messagesMutex.Lock()
	defer m.messagesMutex.Unlock()

	now := time.Now()

	for i := 0; i < len(m.messages); i++ {
		diff := now.Sub(m.messages[i].timestamp)
		if diff < maxMessageAge {
			m.messages = m.messages[i:]
			break
		}
	}
}

func (m *Nuke) nuke(source pkg.User, bot pkg.Sender, channel pkg.Channel, phrase string, scrollbackLength, timeoutDuration time.Duration) {
	if timeoutDuration > 24*time.Hour {
		timeoutDuration = 24 * time.Hour
	}

	matcher := func(msg *nukeMessage) bool {
		return strings.Contains(msg.message.GetText(), phrase)
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

	for i := len(m.messages) - 1; i >= 0; i-- {
		diff := now.Sub(m.messages[i].timestamp)
		if diff > scrollbackLength {
			// We've gone far enough in the buffer, time to exit
			break
		}

		if matcher(&m.messages[i]) {
			targets[m.messages[i].user.GetID()] = m.messages[i].user
		}
	}

	for _, user := range targets {
		bot.Timeout(channel, user, timeoutDurationInSeconds, reason)
	}

	bot.Say(channel, fmt.Sprintf("%s nuked %d users for the phrase %s in the last %s for %s", source.GetName(), len(targets), phrase, scrollbackLength, timeoutDuration))
}

func (m *Nuke) addMessage(channel pkg.Channel, user pkg.User, message pkg.Message) {
	m.messagesMutex.Lock()
	defer m.messagesMutex.Unlock()
	m.messages = append(m.messages, nukeMessage{
		channel:   channel,
		user:      user,
		message:   message,
		timestamp: time.Now(),
	})
}
