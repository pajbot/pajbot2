package commands

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/utils"
)

func NewRaffle() *Raffle {
	return &Raffle{
		participantsUsername: make(map[string]string),
	}
}

type Raffle struct {
	running bool
	points  int64

	// by user ID
	participants         []string
	participantsUsername map[string]string
}

func (c *Raffle) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	cmd := strings.ToLower(parts[0])

	if cmd == "!roffle" {
		if !user.HasChannelPermission(channel, pkg.PermissionRaffle) {
			botChannel.Mention(user, "you do not have the permission to start a raffle")
			return
		}

		if c.running {
			botChannel.Mention(user, "a raffle is already running xd")
			return
		}

		var pointsToRaffle int64
		if len(parts) == 1 {
			pointsToRaffle = 1000
		} else {
			var err error
			pointsToRaffle, err = strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				botChannel.Mention(user, "usage: !raffle 500")
				return
			}
		}

		c.running = true
		c.points = pointsToRaffle

		botChannel.Say("A raffle is now running for " + strconv.FormatInt(pointsToRaffle, 10) + " points PepeS type !join to have a chance to win")

		time.AfterFunc(time.Second*5, func() {
			c.running = false

			// TODO: mutex loooooooooool
			if len(c.participants) == 0 {
				botChannel.Say("no one joined the raffle FeelsBadMan")
				return
			}

			winnerIndex := rand.Intn(len(c.participants))
			winnerID := c.participants[winnerIndex]
			winnerUsername := c.participantsUsername[winnerID]

			var newPoints uint64

			if c.points > 0 {
				_, newPoints = botChannel.Bot().AddPoints(channel, winnerID, uint64(c.points))
			} else {
				newPoints = botChannel.Bot().ForceRemovePoints(channel, winnerID, uint64(utils.Abs64(c.points)))
			}

			botChannel.Say("@" + winnerUsername + ", you won the raffle PogChamp you now have " + strconv.FormatUint(newPoints, 10) + " points")

			c.participants = []string{}
			c.participantsUsername = make(map[string]string)
		})

		// Start raffle, but only if you have permission
		return
	}

	if cmd == "!join" {
		if !c.running {
			// No raffle is running
			return
		}

		if _, ok := c.participantsUsername[user.GetID()]; !ok {
			// User can join the raffle
			c.participantsUsername[user.GetID()] = user.GetName()
			c.participants = append(c.participants, user.GetID())

			botChannel.Mention(user, "you have joined the raffle PepeS")
		}

		return
	}

	botChannel.Mention(user, "how did you get here?")
}
