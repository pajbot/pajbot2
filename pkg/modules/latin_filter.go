package modules

import (
	"encoding/json"
	"log"
	"time"

	twitch "github.com/gempir/go-twitch-irc"
	"github.com/pajlada/pajbot2/filter"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pkg/errors"
)

func maxpenis(a, b int) int {
	if a > b {
		return a
	}

	return b
}

type UnicodeRange struct {
	Start rune
	End   rune
}

type LatinFilter struct {
	server *server

	transparentList  *filter.TransparentList
	unicodeWhitelist []UnicodeRange
}

func NewLatinFilter() *LatinFilter {
	return &LatinFilter{
		server: &_server,

		transparentList: filter.NewTransparentList(),
	}
}

func (m *LatinFilter) addToWhitelist(start, end rune) {
	m.unicodeWhitelist = append(m.unicodeWhitelist, UnicodeRange{start, end})
}

func (m *LatinFilter) Register() error {
	m.transparentList.Add("(/ﾟДﾟ)/")
	m.transparentList.Add("(╯°□°）╯︵ ┻━┻")
	m.transparentList.Add("(╯°Д°）╯︵/(.□ . )")
	m.transparentList.Add("(ノಠ益ಠ)ノ彡┻━┻")
	m.transparentList.Add("୧༼ಠ益ಠ༽୨")
	m.transparentList.Add("༼ ºل͟º ༽")
	m.transparentList.Add("༼つಠ益ಠ༽つ")
	m.transparentList.Add("( ° ͜ʖ͡°)╭∩╮")
	m.transparentList.Add("ᕙ༼ຈل͜ຈ༽ᕗ")
	m.transparentList.Add("ʕ•ᴥ•ʔ")
	m.transparentList.Add("༼▀̿ Ĺ̯▀̿༽")
	m.transparentList.Add("( ͡° ͜🔴 ͡°)")

	err := m.transparentList.Build()
	if err != nil {
		return errors.Wrap(err, "Failed to build transparent list")
	}

	m.addToWhitelist(0x20, 0x7e)       // Basic latin
	m.addToWhitelist(0x1f600, 0x1f64f) // Emojis
	m.addToWhitelist(0x1f300, 0x1f5ff) // "Miscellaneous symbols and pictographs". Includes some emojis like 100
	m.addToWhitelist(0x1f44c, 0x1f44c) // Chatterino?
	m.addToWhitelist(0x206d, 0x206d)   // Chatterino?
	m.addToWhitelist(0x2660, 0x2765)   // Chatterino?

	m.addToWhitelist(0x1f171, 0x1f171) // B emoji
	m.addToWhitelist(0x1f900, 0x1f9ff) // More emojis

	m.addToWhitelist(0x2019, 0x2019) // Scuffed '

	// Rain
	m.addToWhitelist(0x30fd, 0x30fd)
	m.addToWhitelist(0xff40, 0xff40)
	m.addToWhitelist(0x3001, 0x3001)
	m.addToWhitelist(0x2602, 0x2602)

	// From Karl
	m.addToWhitelist(0x1d100, 0x1d1ff)
	m.addToWhitelist(0x1f680, 0x1f6ff)
	m.addToWhitelist(0x2600, 0x26ff)
	m.addToWhitelist(0xfe00, 0xfe0f) // Emoji variation selector 1 to 16
	m.addToWhitelist(0x2012, 0x2015) // Various dashes
	m.addToWhitelist(0x3010, 0x3011) // 【 and 】

	return nil
}

func (m LatinFilter) Name() string {
	return "LatinFilter"
}

func (m LatinFilter) OnMessage(channel string, user pkg.User, message twitch.Message) error {
	if !user.IsModerator() || true {
		lol := struct {
			FullMessage   string
			Message       string
			BadCharacters []rune
			Username      string
			Channel       string
			Timestamp     time.Time
		}{
			FullMessage: message.Text,
			Username:    user.GetName(),
			Channel:     channel,
			Timestamp:   time.Now().UTC(),
		}
		messageRunes := []rune(message.Text)
		transparentStart := time.Now()
		transparentSkipRange := m.transparentList.Find(messageRunes)
		transparentEnd := time.Now()
		if pkg.VerboseBenchmark {
			log.Printf("[% 26s] %s", "TransparentList", transparentEnd.Sub(transparentStart))
		}
		messageLength := len(messageRunes)
		for i := 0; i < messageLength; {
			if skipLength := transparentSkipRange.ShouldSkip(i); skipLength > 0 {
				i = i + skipLength
				continue
			}

			r := messageRunes[i]
			allowed := false

			for _, allowedRange := range m.unicodeWhitelist {
				if r >= allowedRange.Start && r <= allowedRange.End {
					allowed = true
					break
				}
			}

			if !allowed {
				if lol.Message == "" {
					lol.Message = message.Text[maxpenis(0, i-2):len(message.Text)]
				}

				alreadySet := false
				for _, bc := range lol.BadCharacters {
					if bc == r {
						alreadySet = true
						break
					}
				}

				if !alreadySet {
					lol.BadCharacters = append(lol.BadCharacters, r)
				}

			}
			i++
		}

		if lol.Message != "" {
			go func() {
				c := m.server.redis.Pool.Get()
				bytes, _ := json.Marshal(&lol)
				c.Do("LPUSH", "karl_kons", bytes)
				c.Close()
				log.Printf("First bad character: 0x%0x message '%s' from '%s' in '#%s' is disallowed due to our whitelist\n", lol.BadCharacters[0], message.Text, user.GetName(), channel)
			}()
		}
	}

	return nil
}
