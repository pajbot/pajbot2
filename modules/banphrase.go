package modules

import (
	"sync"
	"time"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/filter"
)

/*
Banphrase module settings
*/
type Banphrase struct {
	sync.Mutex
	BannedWords      []filter.BannedWord
	BannedLinks      []filter.BannedLink
	TimeoutDur       [15]int
	Level            map[int]int
	timeoutHistory   map[string]userHistory // username[ban level]
	ResetDuration    time.Duration
	MaxMessageLength int
}

type userHistory struct {
	start   int
	current int
	time    time.Time
}

// these define what was actually wrong with the message
const (
	Link = iota
	BannedWord
	BannedLink
	MessageTooLong
)

// Ensure the module implements the interface properly
var _ Module = (*Banphrase)(nil)

// Init xD
func (module *Banphrase) Init(bot *bot.Bot) {
	// load banned words and links
	module.BannedLinks = []filter.BannedLink{
		filter.BannedLink{
			Link:  "www.com",
			Level: 5,
		},
		filter.BannedLink{
			Link:  "bit.ly",
			Level: 6,
		},
		filter.BannedLink{
			Link:  "google.com",
			Level: 6,
		},
	}
	module.BannedWords = []filter.BannedWord{
		filter.BannedWord{
			Word:  "forsenpuke",
			Level: 2,
		},
		filter.BannedWord{
			Word:  "kappa",
			Level: 1,
		},
		filter.BannedWord{
			Word:  "minik",
			Level: 1,
		},
	}

	// default values
	module.Level = map[int]int{
		Link:           0,
		BannedLink:     0,
		BannedWord:     0,
		MessageTooLong: 0,
	}
	module.TimeoutDur = [15]int{
		0,
		0,
		0,
		0, // send to dashboard
		0, // send to dashboard
		5,
		15,
		30,
		120,
		600,    // 10 min
		1200,   // 20 min
		3600,   // 1 hour
		86400,  // 24h
		604800, // 1 week
		-1,     // perma ban
	}

	module.timeoutHistory = make(map[string]userHistory)
	module.ResetDuration = time.Minute * 1
	module.MaxMessageLength = 250
	go module.gc()
}

// DeInit xD
func (module *Banphrase) DeInit(b *bot.Bot) {

}

// Check xD
func (module *Banphrase) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	if msg.User.Level >= 500 {
		return nil
	}
	lvl, reason := module.runFilters(msg)
	if lvl == 0 {
		return nil
	}
	if lvl == -2 {
		// bot.SendToDashboard xD
		return nil
	}
	l := module.getTimeoutLevel(msg.User.Name, lvl)
	dur := module.TimeoutDur[l]
	if dur < 1 {
		return nil
	}
	b.Timeout(msg.User.Name, dur, reason)
	action.Stop = true

	return nil
}

func (module *Banphrase) getTimeoutLevel(user string, oldLvl int) int {
	var u userHistory
	var ok bool
	reset := true
	if u, ok = module.timeoutHistory[user]; ok {
		reset = time.Now().After(u.time.Add(module.ResetDuration))
		log.Debug(u)
		log.Debug(oldLvl)
		if u.current+1 < oldLvl {
			reset = true
		} else if u.start+3 > u.current {
			// max 3 levels over the start level
			u.current++
			u.time = time.Now()
		}
	}
	if reset {
		log.Debug("RESET LEVEL")
		u = userHistory{
			start:   oldLvl,
			current: oldLvl,
			time:    time.Now(),
		}
	}
	module.Lock()
	module.timeoutHistory[user] = u
	module.Unlock()
	return u.current
}

func (module *Banphrase) runFilters(msg *common.Msg) (int, string) {
	var lvl int
	var reason string
	m := msg.Text // TODO: confusables, github.com/FiloSottile/tr39-confusables this is kinda shitty
	links := filter.LinkFilter(m)
	//log.Debug(links)
	if msg.User.Level < 250 && len(links) > 0 {
		lvl = module.Level[Link]
		reason = "matched link filter"
	}

	if l := filter.ContainsLink(links, module.BannedLinks); l > lvl {
		lvl = l
		reason = "contains banned link"
	}
	if l := filter.ContainsWord(m, module.BannedWords); l > lvl {
		lvl = l
		reason = "contains banned word"
	}
	if filter.MessageLength(msg) > module.MaxMessageLength {
		lvl = module.Level[MessageTooLong]
		reason = "message too long"
	}
	return lvl, reason
}

func (module *Banphrase) gc() {
	for {
		time.Sleep(module.ResetDuration * 2)

		for user, data := range module.timeoutHistory {
			if time.Now().After(data.time.Add(module.ResetDuration)) {
				module.Lock()
				delete(module.timeoutHistory, user)
				log.Debug("removed user: ", user)
				module.Unlock()
			}
		}
	}
}
