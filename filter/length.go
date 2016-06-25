package filter

import (
	"strings"
	"unicode"

	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/plog"
)

var log = plog.GetLogger()

// Length xD should
// should this be a filter or a module?
type Length struct {
}

var _ Filter = (*Length)(nil)

// Run filter
func (length *Length) Run(_ string, msg *common.Msg, action *BanAction) {
	var msgLength int
	var emoteLength int
	m := msg.Text
	for _, emote := range msg.Emotes {
		emoteLength = (emote.SizeX * emote.SizeY) / 100
		// 7 for twitch emotes, 16 for NaM
		if emote.IsGif {
			emoteLength = emoteLength * 2
		}
		emoteLength = emoteLength * emote.Count
		m = strings.Replace(m, emote.Name, "", -1)
		// TODO: parse emote names
		// log.Debugf("Emote: %s [Length: %d]", emote.Name, emoteLength)
		msgLength += emoteLength
	}
	runes := []rune(m)
	for _, l := range runes {
		if unicode.IsLetter(l) || unicode.IsDigit(l) || unicode.IsSpace(l) || unicode.IsPunct(l) {
			msgLength++
		} else if unicode.IsGraphic(l) {
			msgLength += 10
		} else {
			msgLength += 2
		}
	}
	msg.Length = msgLength
	log.Debugf("Message: %s [Length: %d]", msg.Text, msg.Length)
}
