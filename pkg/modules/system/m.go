package system

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/modules"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

const id = "system"
const name = "System"

func init() {
	modules.Register(id, func() pkg.ModuleSpec {
		return modules.NewSpec(id, name, true, newModule)
	})
}

type module struct {
	mbase.Base
}

func newModule(b *mbase.Base) pkg.Module {
	m := &module{
		Base: *b,
	}

	return m
}

func (m *module) OnEventSubNotification(event pkg.EventSubNotificationEvent) pkg.Actions {
	switch event.Notification.Subscription.Type {
	case helix.EventSubTypeChannelFollow:
		var followEvent helix.EventSubChannelFollowEvent
		err := json.NewDecoder(bytes.NewReader(event.Notification.Event)).Decode(&followEvent)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return twitchactions.Sayf("%s just followed the stream xD", followEvent.UserLogin)
	default:
		fmt.Println("Got unhandled eventsub notification:", event)
	}
	return nil
}
