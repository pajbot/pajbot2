package common

import (
	"bytes"
	"fmt"

	"github.com/urakozz/go-emoji"
)

var parser = emoji.NewEmojiParser()

// ParseEmojis prases emojis from the message Text
func ParseEmojis(msg *Msg) {
	emoteCount := make(map[string]*Emote)
	_ = parser.ReplaceAllStringFunc(msg.Text, func(s string) string {
		byteArray := []byte(s)
		log.Debug(byteArray)
		log.Debug(bytes.Runes(byteArray))
		if emote, ok := emoteCount[s]; ok {
			emote.Count++
		} else {
			emoteCount[s] = &Emote{
				ID:    fmt.Sprintf("%x", bytes.Runes(byteArray)[0]),
				Type:  "emoji",
				SizeX: 28,
				SizeY: 28,
				IsGif: false,
				Count: 1,
			}
		}
		return ""
	})

	for _, emote := range emoteCount {
		msg.Emotes = append(msg.Emotes, *emote)
	}
}
