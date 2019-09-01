package modules

import (
	"math/rand"
	"strings"

	"github.com/pajbot/pajbot2/pkg"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
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
	mbase.Base

	state string

	entrants []string
}

func newGiveaway(b mbase.Base) pkg.Module {
	return &giveaway{
		Base: b,

		state: "inactive",
	}
}

const forsen25ID = "300378550"

func (m *giveaway) OnMessage(event pkg.MessageEvent) pkg.Actions {
	const giveawayEmote = forsen25ID

	message := event.Message
	user := event.User

	text := message.GetText()

	// Commands
	if user.IsModerator() {
		if strings.HasPrefix(text, "!25start") {
			if m.state == "inactive" {
				m.state = "started"
				m.entrants = []string{}
				return twitchactions.Say("Started giveaway")
			}

			return twitchactions.Say("Giveaway already started")
		}

		if strings.HasPrefix(text, "!25stop") {
			if m.state == "started" {
				m.state = "inactive"
				return twitchactions.Say("Stopped accepting people into the giveaway")
			}
		}

		if strings.HasPrefix(text, "!25draw") {
			if len(m.entrants) == 0 {
				return twitchactions.Say("No one has joined the giveaway")
			}

			winnerIndex := rand.Intn(len(m.entrants))
			winnerUsername := m.entrants[winnerIndex]

			m.entrants = append(m.entrants[:winnerIndex], m.entrants[winnerIndex+1:]...)

			return twitchactions.Say(winnerUsername + " just won the sub emote giveaway PogChamp")
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

			return twitchactions.Mention(event.User, "you have been entered into the sub emote giveaway")
		}
	}

	return nil
}
