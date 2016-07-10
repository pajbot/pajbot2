package modules

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/command"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/helper"
)

type activeBingo struct {
}

func (bingo *activeBingo) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) {

}

/*
Bingo xD
*/
type Bingo struct {
	commandHandler command.Handler
	activeBingos   map[string]*activeBingo
}

// Ensure the module implements the interface properly
var _ Module = (*Bingo)(nil)

func (module *Bingo) usageCommand(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	b.Say("Usage: !bango type reward. Valid types: twitch/bttv/number")
}

var bingoRunning = false
var bingoCancelChannel = make(chan bool)
var bingoMessageChannel = make(chan *common.Msg)

func (module *Bingo) bingoNumber(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	const usageString = "Usage: !bango number 1-1000 500"
	var numLow int
	var numHigh int
	arguments := helper.GetTriggersN(msg.Text, 2)

	if bingoRunning {
		b.Say("A bingo is already running. Use !bango cancel to cancel the current bingo")
		return
	}

	if len(arguments) < 2 {
		b.Say(usageString)
		return
	}

	lowHigh := strings.Split(arguments[0], "-")
	if len(lowHigh) != 2 {
		b.Say(usageString)
		return
	}

	numLow, err := strconv.Atoi(lowHigh[0])
	if err != nil {
		log.Errorf("Error in bingoNumbeR: %s", err)
		return
	}

	numHigh, err = strconv.Atoi(lowHigh[1])
	if err != nil {
		log.Errorf("Error in bingoNumbeR: %s", err)
		return
	}

	if numLow < 0 {
		b.Say("Lowest number must be higher or equal to 0")
		return
	}

	if numLow >= numHigh {
		b.Say("First number must be less than second number")
		return
	}

	pointReward, err := strconv.Atoi(arguments[1])
	if err != nil {
		log.Errorf("Error convinerting pointReward: %s", err)
		return
	}

	bingoRunning = true

	b.Sayf("say a number between %d and %d", numLow, numHigh)

	go func(numLow int, numHigh int, pointReward int) {
		defer func() {
			bingoRunning = false
		}()

		// Select a random number
		winningNumber := numLow + rand.Intn(numHigh-numLow)

		for {
			select {
			case _ = <-bingoCancelChannel:
				return
			case newMessage := <-bingoMessageChannel:
				// Check if the message is good!

				// the message must have the number as the first "split"
				msgSplit := strings.Split(newMessage.Text, " ")
				potentialNumber, err := strconv.Atoi(msgSplit[0])
				if err == nil {
					if potentialNumber == winningNumber {
						// The user guessed right.
						// WHAT DO WE DO?
						b.Sayf("%s just won the number bingo with the guess %d! He wins %d points",
							newMessage.User.DisplayName, winningNumber, pointReward)
						b.Redis.IncrPoints(b.Channel.Name, newMessage.User.Name, pointReward)
						bingoRunning = false
						return
					}
				}
			}
		}
	}(numLow, numHigh, pointReward)
}

func (module *Bingo) topSpammerOnline(b *bot.Bot, msg *common.Msg, action *bot.Action) {
}

func (module *Bingo) topSpammerOffline(b *bot.Bot, msg *common.Msg, action *bot.Action) {
}

func (module *Bingo) topSpammerTotal(b *bot.Bot, msg *common.Msg, action *bot.Action) {
}

// Init xD
func (module *Bingo) Init(bot *bot.Bot) {
	numberCommand := &command.FuncCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"number",
				"digit",
				"num",
			},
			Level: 500,
		},
		Function: module.bingoNumber,
	}
	usageCommand := &command.FuncCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"usage",
				"help",
			},
		},
		Function: module.usageCommand,
	}
	bingoCommand := &command.NestedCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"bango",
			},
			Level: 500,
		},
		Commands: []command.Command{
			usageCommand,
			numberCommand,
		},
		DefaultCommand:  usageCommand,
		FallbackCommand: usageCommand,
	}
	module.commandHandler.AddCommand(bingoCommand)
}

// Check xD
func (module *Bingo) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	if bingoRunning {
		bingoMessageChannel <- msg
	}

	return module.commandHandler.Check(b, msg, action)
}
