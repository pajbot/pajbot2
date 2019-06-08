package modules

import (
	"math/rand"
	"strings"

	"github.com/pajbot/pajbot2/pkg"
)

func init() {
	Register("giveaway", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:    "giveaway",
			name:  "Giveaway",
			maker: newGiveaway,
		}
	})
}

type giveaway struct {
	base

	state string

	entrants []string
}

func newGiveaway(b base) pkg.Module {
	return &giveaway{
		base: b,

		state: "inactive",
	}
}

const forsen25ID = "300378550"

func (m *giveaway) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	const giveawayEmote = forsen25ID

	text := message.GetText()

	// Commands
	if user.IsModerator() || user.IsBroadcaster(bot.Channel()) {
		if strings.HasPrefix(text, "!25start") {
			if m.state == "inactive" {
				m.state = "started"
				m.entrants = []string{}
				bot.Say("Started giveaway")

				return nil
			}

			bot.Say("Giveaway already started")
			return nil
		}

		if strings.HasPrefix(text, "!25stop") {
			if m.state == "started" {
				m.state = "inactive"
				bot.Say("Stopped accepting people into the giveaway")

				return nil
			}
		}

		if strings.HasPrefix(text, "!25draw") {
			if len(m.entrants) == 0 {
				bot.Say("No one has joined the giveaway")
				return nil
			}

			winnerIndex := rand.Intn(len(m.entrants))
			winnerUsername := m.entrants[winnerIndex]
			bot.Say(winnerUsername + " just won the sub emote giveaway PogChamp")

			m.entrants = append(m.entrants[:winnerIndex], m.entrants[winnerIndex+1:]...)

			return nil
		}
	}

	if m.state == "started" {
		enterGiveaway := false

		reader := message.GetTwitchReader()
		for reader.Next() {
			emote := reader.Get()
			if emote.GetID() == giveawayEmote {
				// bot.Say(fmt.Sprintf("%#v", emote))
				// User wants to enter the giveaway
				enterGiveaway = true
				break
			}
		}

		if enterGiveaway {
			for _, entrant := range m.entrants {
				if entrant == user.GetName() {
					// User has already joined
					return nil
				}

			}
			m.entrants = append(m.entrants, user.GetName())

			bot.Mention(user, "you have been entered into the sub emote giveaway")
			return nil
		}
	}

	return nil
}
