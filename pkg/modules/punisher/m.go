package punisher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/modules"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
)

const id = "punisher"
const name = "Punisher"

func init() {
	modules.Register(id, func() pkg.ModuleSpec {
		return modules.NewSpec(id, name, false, newModule)
	})
}

type punishment struct {
	timeout timeout

	username string
}

const (
	lookbehindDuration = 5 * time.Second
	punishmentInterval = 2 * time.Second
)

type module struct {
	mbase.Base

	timeouts map[string]*timeout

	scheduledPunishments      map[string]punishment
	scheduledPunishmentsMutex sync.Mutex

	ticker *time.Ticker

	cancel     context.Context
	cancelFunc context.CancelFunc
}

func (m *module) handleClearChat(args map[string]interface{}) error {
	iMessage, ok := args["message"]
	if !ok {
		return nil
	}

	message, ok := iMessage.(*twitch.ClearChatMessage)
	if !ok {
		return nil
	}

	newTimeout := &timeout{
		received: message.Time,
		end:      message.Time.Add(time.Duration(message.BanDuration) * time.Second),

		duration: message.BanDuration,
	}

	oldTimeout, isOldTimeout := m.timeouts[message.TargetUserID]
	if !isOldTimeout {
		m.timeouts[message.TargetUserID] = newTimeout
		return nil
	}

	breakoff := oldTimeout.received.Add(lookbehindDuration)

	if message.Time.After(breakoff) {
		oldTimeout.end = newTimeout.end
		oldTimeout.duration = newTimeout.duration
		oldTimeout.received = newTimeout.received
		return nil
	}

	m.scheduledPunishmentsMutex.Lock()
	defer m.scheduledPunishmentsMutex.Unlock()
	if newTimeout.IsSmaller(oldTimeout, 10) {
		fmt.Println("[PUNISHER] Schedule re-doing timeout on", message.TargetUsername)
		m.scheduledPunishments[message.TargetUserID] = punishment{
			timeout:  *oldTimeout,
			username: message.TargetUsername,
		}
	} else {
		oldTimeout.duration = newTimeout.duration
		oldTimeout.end = newTimeout.end

		fmt.Println("[PUNISHER] Removing scheduled timeout on", message.TargetUsername)
		delete(m.scheduledPunishments, message.TargetUserID)
	}
	oldTimeout.received = message.Time

	return nil
}

func (m *module) startPunisher() {
	go func() {
		for {
			select {
			case <-m.ticker.C:
				m.scheduledPunishmentsMutex.Lock()

				for _, punishment := range m.scheduledPunishments {
					if punishment.timeout.duration == 0 {
						fmt.Println("[PUNISHER] Perform ban punishment on user:", punishment.username)
						m.BotChannel().Ban(m.BotChannel().Bot().MakeUser(punishment.username), "PUNISHER")
					} else {
						newDuration := int(punishment.timeout.Seconds())
						fmt.Println("[PUNISHER] Perform timeout punishment on user:", punishment.username, newDuration)
						m.BotChannel().SingleTimeout(m.BotChannel().Bot().MakeUser(punishment.username), newDuration, "PUNISHER")
					}
				}

				m.scheduledPunishments = map[string]punishment{}

				m.scheduledPunishmentsMutex.Unlock()

			case <-m.cancel.Done():
				fmt.Println("CANCEL")
				return
			}
		}
	}()
}

func (m *module) Disable() error {
	// Stop the ticker from firing
	m.ticker.Stop()

	// Cancel the ticker selector
	m.cancelFunc()

	return m.Base.Disable()
}

func newModule(b *mbase.Base) pkg.Module {
	m := &module{
		Base: *b,

		timeouts:             map[string]*timeout{},
		scheduledPunishments: map[string]punishment{},

		ticker: time.NewTicker(punishmentInterval),
	}

	m.cancel, m.cancelFunc = context.WithCancel(context.Background())

	err := m.Listen("on_clearchat", m.handleClearChat, 100)
	if err != nil {
		fmt.Println("ERROR LISTENING XD", err)
	}

	m.startPunisher()

	return m
}
