package modules

import (
	"math/rand"
	"strings"

	"github.com/pajlada/pajbot2/pkg"
)

type Giveaway struct {
	server *server

	state string

	entrants []string

	Sender pkg.Sender
}

func NewGiveaway(sender pkg.Sender) *Giveaway {
	return &Giveaway{
		server: &_server,
		state:  "inactive",
		Sender: sender,
	}
}

func (m Giveaway) Name() string {
	return "Giveaway"
}

func (m Giveaway) Register() error {
	return nil
}

const forsen25ID = "1361602"
const pajlada25ID = "908917"
const pajaWID = "80481"

func (m Giveaway) OnWhisper(bot pkg.Sender, user pkg.User, message pkg.Message) error {
	return nil
}

func (m *Giveaway) OnMessage(bot pkg.Sender, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	// if channel.GetChannel() != "forsen" {
	// 	return nil
	// }

	const giveawayEmote = forsen25ID

	text := message.GetText()

	// Commands
	if user.IsModerator() || user.GetName() == "forsen" || user.GetName() == "pajlada" {
		if strings.HasPrefix(text, "!25start") {
			if m.state == "inactive" {
				m.state = "started"
				m.entrants = []string{}
				m.Sender.Say(channel, "Started giveaway")

				return nil
			} else {
				m.Sender.Say(channel, "Giveaway already started")
				return nil
			}
		}

		if strings.HasPrefix(text, "!25stop") {
			if m.state == "started" {
				m.state = "inactive"
				m.Sender.Say(channel, "Stopped accepting people into the giveaway")

				return nil
			}
		}

		if strings.HasPrefix(text, "!25draw") {
			if len(m.entrants) == 0 {
				m.Sender.Say(channel, "No one has joined the giveaway")
				return nil
			}

			winnerIndex := rand.Intn(len(m.entrants))
			winnerUsername := m.entrants[winnerIndex]
			m.Sender.Say(channel, winnerUsername+" just won the sub emote giveaway PogChamp")

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
				// m.Sender.Say(channel, fmt.Sprintf("%#v", emote))
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

			m.Sender.Say(channel, "@"+user.GetName()+", you have been entered into the sub emote giveaway")
			return nil
		}
	}

	return nil
}
