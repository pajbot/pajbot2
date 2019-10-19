package modules

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

type giveawayState int

const (
	giveawayStateStopped giveawayState = iota
	giveawayStateStarted
)

func init() {
	Register("giveaway", func() pkg.ModuleSpec {
		return &Spec{
			id:    "giveaway",
			name:  "Giveaway",
			maker: newGiveaway,

			parameters: map[string]pkg.ModuleParameterSpec{
				"EmoteID": func() pkg.ModuleParameter {
					return newStringParameter(parameterSpec{
						Description: "Emote ID that needs to exist in a message for a user to join the giveaway",
					})
				},
				"EmoteName": func() pkg.ModuleParameter {
					return newStringParameter(parameterSpec{
						Description: "Name of the emote that needs to exist in a message for a user to join the giveaway",
					})
				},
			},
		}
	})
}

type giveaway struct {
	mbase.Base

	commands pkg.CommandsManager

	state giveawayState

	entrants []string

	emoteID   string
	emoteName string
}

func newGiveaway(b mbase.Base) pkg.Module {
	m := &giveaway{
		Base: b,

		state: giveawayStateStopped,

		commands: commands.NewCommands(),
	}

	m.Parameters()["EmoteID"].Link(&m.emoteID)
	m.Parameters()["EmoteName"].Link(&m.emoteName)

	m.commands.Register([]string{"!25start"}, giveawayCmdStop{m: m})
	m.commands.Register([]string{"!25stop"}, giveawayCmdStop{m: m})
	m.commands.Register([]string{"!25draw"}, giveawayCmdDraw{m: m})
	m.commands.Register([]string{"!25config"}, giveawayCmdConfig{m: m})

	return m
}

func (m *giveaway) start() {
	m.state = giveawayStateStarted
	m.entrants = []string{}
}

func (m *giveaway) started() bool {
	return m.state == giveawayStateStarted
}

func (m *giveaway) stop() {
	m.state = giveawayStateStopped
}

func (m *giveaway) stopped() bool {
	return m.state == giveawayStateStopped
}

func (m *giveaway) OnMessage(event pkg.MessageEvent) pkg.Actions {
	if actions := m.commands.OnMessage(event); actions != nil {
		return actions
	}

	giveawayEmote := m.emoteID
	message := event.Message
	user := event.User

	if m.started() {
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
