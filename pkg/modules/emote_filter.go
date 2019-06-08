package modules

import (
	"fmt"
	"strings"

	"github.com/pajbot/pajbot2/pkg"
)

func init() {
	Register("emote_filter", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:    "emote_limit",
			name:  "Emote limit",
			maker: newEmoteFilter,
		}
	})
}

type limitConsequence struct {
	limit int

	baseDuration int

	extraDuration int
}

type emoteFilter struct {
	base

	emoteLimits    map[string]limitConsequence
	combinedLimits int
}

func newEmoteFilter(b base) pkg.Module {
	m := &emoteFilter{
		base: b,

		emoteLimits: make(map[string]limitConsequence),

		combinedLimits: 4,
	}

	// FIXME
	m.Initialize()

	return m
}

func (m *emoteFilter) Initialize() {
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
	m.emoteLimits["bttvNice"] = limitConsequence{
		limit:         3,
		baseDuration:  300,
		extraDuration: 50,
	}
}

func (m *emoteFilter) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
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
		action.Set(pkg.Timeout{timeoutDuration, "Don't overuse " + strings.Join(overusedEmotes, ", ")})
	} else if combinedLimits > m.combinedLimits {
		action.Set(pkg.Timeout{combinedLimits * 120, "Don't overuse big emotes"})
	}

	return nil
}
