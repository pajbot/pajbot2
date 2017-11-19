package bot

import "github.com/pajlada/pajbot2/common"

// LoadBttvEmotes should load emotes from redis, but this should do for now
func (bot *Bot) LoadBttvEmotes() {
	// this loads bttv channel emotes
	/*
		apirequest.BTTV.GetChannel(bot.Channel.Name,
			func(channel gobttv.ChannelResponse) {
				bot.Channel.Emotes.Bttv = make(map[string]common.Emote)
				bot.Channel.Emotes.BttvLastUpdate = time.Now()

				for _, emote := range channel.Emotes {
					bot.Channel.Emotes.Bttv[emote.Code] = ParseBTTVChannelEmote(emote)
				}
			},
			func(statusCode int, statusMessage, errorMessage string) {
				// We ignore 404 errors, it just means he doesn't have a BTTV account
				if statusCode != 404 {
					log.Printf("Error fetching Channel BTTV Emotes (%s)", bot.Channel.Name)
					log.Printf("Status code: %d", statusCode)
					log.Printf("Status message: %s", statusMessage)
					log.Printf("Error message: %s", errorMessage)
				}
			}, func(err error) {
				log.Printf("Internal error: %s", err)
			})
	*/
}

// regex would probably be better but im a regex noob ¯\_(ツ)_/¯
func (bot *Bot) parseEmotes(msg *common.Msg) {
	/*
		m := strings.Split(msg.Text, " ")
		emoteCount := make(map[string]*common.Emote)
		for _, word := range m {
			if emote, ok := emoteCount[word]; ok {
				emote.Count++
			} else if emote, ok := bot.Channel.Emotes.Bttv[word]; ok {
				emoteCount[word] = &emote
			} else if emote, ok := GlobalEmotes.Bttv[word]; ok {
				emoteCount[word] = &emote
			} else if emote, ok := bot.Channel.Emotes.FrankerFaceZ[word]; ok {
				emoteCount[word] = &emote
			} else if emote, ok := GlobalEmotes.FrankerFaceZ[word]; ok {
				emoteCount[word] = &emote
			}
		}

		for _, emote := range emoteCount {
			msg.Emotes = append(msg.Emotes, *emote)
		}
	*/
}

// LoadFFZEmotes should load emotes from redis, but this should do for now
func (bot *Bot) LoadFFZEmotes() {
	// this loads channel emotes
	/*
		apirequest.FFZ.GetRoom(bot.Channel.Name,
			func(room goffz.RoomResponse) {
				bot.Channel.Emotes.FrankerFaceZ = make(map[string]common.Emote)
				bot.Channel.Emotes.FrankerFaceZLastUpdate = time.Now()

				for _, set := range room.Sets {
					for _, emote := range set.Emoticons {
						bot.Channel.Emotes.FrankerFaceZ[emote.Name] = ParseFrankerFaceZEmote(emote)
					}
				}
			},
			func(statusCode int, statusMessage, errorMessage string) {
				// We ignore 404 errors, it just means he doesn't have a FFZ account
				if statusCode != 404 {
					log.Printf("Error fetching Channel FFZ Emotes (%s)", bot.Channel.Name)
					log.Printf("Status code: %d", statusCode)
					log.Printf("Status message: %s", statusMessage)
					log.Printf("Error message: %s", errorMessage)
				}
			}, func(err error) {
				log.Printf("Internal error: %s", err)
			})
	*/
}
