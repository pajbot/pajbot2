package modules

import (
	"fmt"
	"time"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/modules/datastructures"
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

type latinFilter struct {
	botChannel pkg.BotChannel

	server *server

	transparentList  *datastructures.TransparentList
	unicodeWhitelist []UnicodeRange
}

func newLatinFilter() pkg.Module {
	return &latinFilter{
		server: &_server,

		transparentList: datastructures.NewTransparentList(),
	}
}

var latinFilterSpec = moduleSpec{
	id:    "latin_filter",
	name:  "Latin filter",
	maker: newLatinFilter,
}

func (m *latinFilter) addToWhitelist(start, end rune) {
	m.unicodeWhitelist = append(m.unicodeWhitelist, UnicodeRange{start, end})
}

func (m *latinFilter) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

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
	m.addToWhitelist(0xb0, 0xb0)     // degrees symbol

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

func (m *latinFilter) Disable() error {
	return nil
}

func (m *latinFilter) Spec() pkg.ModuleSpec {
	return &latinFilterSpec
}

func (m *latinFilter) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m *latinFilter) OnWhisper(bot pkg.BotChannel, source pkg.User, message pkg.Message) error {
	return nil
}

func (m *latinFilter) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	if !user.IsModerator() || true {
		text := message.GetText()

		lol := struct {
			FullMessage   string
			Message       string
			BadCharacters []rune
			Username      string
			Channel       string
			Timestamp     time.Time
		}{
			FullMessage: text,
			Username:    user.GetName(),
			Channel:     bot.Channel().GetName(),
			Timestamp:   time.Now().UTC(),
		}
		messageRunes := []rune(text)
		transparentStart := time.Now()
		transparentSkipRange := m.transparentList.Find(messageRunes)
		transparentEnd := time.Now()
		if pkg.VerboseBenchmark {
			fmt.Printf("[% 26s] %s", "TransparentList", transparentEnd.Sub(transparentStart))
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
					lol.Message = text[maxpenis(0, i-2):len(text)]
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
	}

	return nil
}
