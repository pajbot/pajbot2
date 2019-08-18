package commands

import (
	"github.com/pajbot/pajbot2/pkg"
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

func (c *Raffle) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	// FIXME: Re-implement (POINTS)
	// cmd := strings.ToLower(parts[0])

	// if cmd == "!roffle" {
	// 	if !event.User.HasChannelPermission(event.Channel, pkg.PermissionRaffle) {
	// 		return twitchactions.Mention(event.User, "you do not have the permission to start a raffle")
	// 	}

	// 	if c.running {
	// 		return twitchactions.Mention(event.User, "a raffle is already running xd")
	// 	}

	// 	var pointsToRaffle int64
	// 	if len(parts) == 1 {
	// 		pointsToRaffle = 1000
	// 	} else {
	// 		var err error
	// 		pointsToRaffle, err = strconv.ParseInt(parts[1], 10, 32)
	// 		if err != nil {
	// 			return twitchactions.Mention(event.User, "usage: !raffle 500")
	// 		}
	// 	}

	// 	c.running = true
	// 	c.points = pointsToRaffle

	// 	botChannel.Say("A raffle is now running for " + strconv.FormatInt(pointsToRaffle, 10) + " points PepeS type !join to have a chance to win")

	// 	time.AfterFunc(time.Second*5, func() {
	// 		c.running = false

	// 		// TODO: mutex loooooooooool
	// 		if len(c.participants) == 0 {
	// 			return twitchactions.Say("no one joined the raffle FeelsBadMan")
	// 		}

	// 		winnerIndex := rand.Intn(len(c.participants))
	// 		winnerID := c.participants[winnerIndex]
	// 		winnerUsername := c.participantsUsername[winnerID]

	// 		var newPoints uint64

	// 		if c.points > 0 {
	// 			_, newPoints = botChannel.Bot().AddPoints(channel, winnerID, uint64(c.points))
	// 		} else {
	// 			newPoints = botChannel.Bot().ForceRemovePoints(channel, winnerID, uint64(utils.Abs64(c.points)))
	// 		}

	// 		botChannel.Say("@" + winnerUsername + ", you won the raffle PogChamp you now have " + strconv.FormatUint(newPoints, 10) + " points")

	// 		c.participants = []string{}
	// 		c.participantsUsername = make(map[string]string)
	// 	})

	// 	// Start raffle, but only if you have permission
	// 	return
	// }

	// if cmd == "!join" {
	// 	if !c.running {
	// 		// No raffle is running
	// 		return
	// 	}

	// 	if _, ok := c.participantsUsername[event.User.GetID()]; !ok {
	// 		// User can join the raffle
	// 		c.participantsUsername[event.User.GetID()] = event.User.GetName()
	// 		c.participants = append(c.participants, event.User.GetID())

	// 		botChannel.Mention(event.User, "you have joined the raffle PepeS")
	// 	}

	// 	return
	// }

	// botChannel.Mention(event.User, "how did you get here")

	return nil
}
