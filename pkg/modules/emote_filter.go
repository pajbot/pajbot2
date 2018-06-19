package modules

import (
	"fmt"
	"strings"

	"github.com/pajlada/pajbot2/pkg"
)

type limitConsequence struct {
	limit int

	baseDuration int

	extraDuration int
}

type EmoteFilter struct {
	server *server

	emoteLimits    map[string]limitConsequence
	combinedLimits int
}

func NewEmoteFilter() *EmoteFilter {
	return &EmoteFilter{
		server: &_server,

		emoteLimits: make(map[string]limitConsequence),
	}
}

func (m *EmoteFilter) Register() error {
	m.emoteLimits["NaM"] = limitConsequence{
		limit:         2,
		baseDuration:  300,
		extraDuration: 60,
	}
	m.emoteLimits["SexPanda"] = limitConsequence{
		limit:         2,
		baseDuration:  300,
		extraDuration: 60,
	}
	m.emoteLimits["TaxiBro"] = limitConsequence{
		limit:         2,
		baseDuration:  300,
		extraDuration: 60,
	}
	m.emoteLimits["FishMoley"] = limitConsequence{
		limit:         2,
		baseDuration:  300,
		extraDuration: 60,
	}
	m.emoteLimits["YetiZ"] = limitConsequence{
		limit:         2,
		baseDuration:  300,
		extraDuration: 60,
	}
	m.combinedLimits = 4
	return nil
}

func (m EmoteFilter) Name() string {
	return "EmoteFilter"
}

func (m EmoteFilter) OnWhisper(source pkg.User, message pkg.Message) error {
	return nil
}

func (m EmoteFilter) OnMessage(source pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	// BTTV Emotes
	reader := message.GetBTTVReader()
	timeoutDuration := 0
	overusedEmotes := []string{}
	combinedLimits := 0
	for reader.Next() {
		emote := reader.Get()

		if limit, ok := m.emoteLimits[emote.GetName()]; ok {
			if emote.GetCount() > limit.limit {
				timeoutDuration += limit.baseDuration
				timeoutDuration += (emote.GetCount() - limit.limit - 1) * limit.extraDuration
				overusedEmotes = append(overusedEmotes, fmt.Sprintf("%s(%d)", emote.GetName(), emote.GetCount()))
			} else {
				combinedLimits += emote.GetCount()
			}
		}
	}

	if timeoutDuration > 0 {
		action.SetTimeout(timeoutDuration, "Don't overuse "+strings.Join(overusedEmotes, ", "))
	} else if combinedLimits > m.combinedLimits {
		action.SetTimeout(combinedLimits*120, "Don't overuse big emotes")
	}

	return nil
}
