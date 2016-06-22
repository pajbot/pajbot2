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
	var msglength int
	m := msg.Message
	for _, emote := range msg.Emotes {
		var elen int
		elen = (emote.SizeX * emote.SizeY) / 100
		// 7 for twitch emotes, 16 for NaM
		if emote.IsGif {
			elen = elen * 2
		}
		elen = elen * emote.Count
		m = strings.Replace(m, emote.Name, "", -1)
		// TODO: parse emote names
		log.Debug(elen)
		msglength += elen
	}
	runes := []rune(m)
	for _, l := range runes {
		if unicode.IsLetter(l) || unicode.IsDigit(l) || unicode.IsSpace(l) || unicode.IsPunct(l) {
			msglength++
		} else if unicode.IsGraphic(l) {
			msglength += 10
		} else {
			msglength += 2
		}
	}
	msg.Length = msglength
	log.Debug(msg.Length)
}
