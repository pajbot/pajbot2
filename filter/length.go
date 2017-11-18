package filter

import (
	"strings"
	"unicode"

	"github.com/pajlada/pajbot2/common"
)

// MessageLength , should this be a module?
// if not we need filter settings
func MessageLength(msg *common.Msg) int {
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
	return msgLength
}
