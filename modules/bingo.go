package modules

import (
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

func (module *Bingo) bingoCancel(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	if !bingoRunning {
		b.Say("No bingo is running...")
		return
	}

	// cancel bingo
	bingoCancelChannel <- true
}

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

	// Select a random number
	winningNumber, err := helper.RandIntN(numLow, numHigh)
	if err != nil {
		b.Sayf("Invalid numbers: %s", err)
		return
	}

	bingoRunning = true

	b.Sayf("A number bingo has been started! To win the %d point reward, guess the right number between %d and %d (both inclusive)", pointReward, numLow, numHigh)

	go func(winningNumber int, pointReward int) {
		defer func() {
			bingoRunning = false
		}()

		for {
			select {
			case _ = <-bingoCancelChannel:
				b.Say("The number bingo has been cancelled")
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
	}(winningNumber, pointReward)
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
	cancelCommand := &command.FuncCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"cancel",
				"stop",
			},
			Level: 500,
		},
		Function: module.bingoCancel,
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
			cancelCommand,
		},
		DefaultCommand:  usageCommand,
		FallbackCommand: usageCommand,
	}
	module.commandHandler.AddCommand(bingoCommand)
}

// DeInit xD
func (module *Bingo) DeInit(b *bot.Bot) {

}

// Check xD
func (module *Bingo) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	if bingoRunning {
		bingoMessageChannel <- msg
	}

	return module.commandHandler.Check(b, msg, action)
}
