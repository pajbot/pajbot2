package moderation

import (
	"fmt"
	"log"
	"net/http"

	twitch "github.com/gempir/go-twitch-irc/v2"
	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg"
	pb2twitch "github.com/pajbot/pajbot2/pkg/twitch"
	"github.com/pajbot/pajbot2/pkg/users"
	"github.com/pajbot/pajbot2/pkg/web/state"
	"github.com/pajbot/utils"
	_ "github.com/swaggo/echo-swagger" // echo-swagger middleware
)

func apiCheckMessageMissingVariables(w http.ResponseWriter, r *http.Request) {
	utils.WebWriteError(w, 400, "Missing message parameter")
}

type filterData struct {
	MuteType pkg.MuteType `json:"mute_type"`
	Reason   string       `json:"reason"`
}

type CheckMessageSuccessResponse struct {
	Banned       bool         `json:"banned"`
	InputMessage string       `json:"input_message"`
	FilterData   []filterData `json:"filter_data,omitempty"`
}

// apiCheckMessage godoc
// @Summary Check a message in a bots filter
// @Produce  json
// @Param channelID path string true "ID of channel to run the test in"
// @Param message query string true "message to test against the bots filters"
// @Success 200 {object} moderation.CheckMessageSuccessResponse
// @Failure 404 {object} utils.WebAPIError
// @Router /api/channel/{channelID}/moderation/check_message [get]
func apiCheckMessage(w http.ResponseWriter, r *http.Request) {
	c := state.Context(w, r)

	vars := mux.Vars(r)

	message := vars["message"]

	if c.Channel == nil {
		// ERROR: Missing required ChannelID parameter
		utils.WebWriteError(w, 400, "bad channel")
		return
	}

	if message == "" {
		// ERROR: Missing required ChannelID parameter
		utils.WebWriteError(w, 400, "Missing required message parameter")
		return
	}

	fmt.Println("Checking message:", message)

	// if !webutils.RequirePermission(w, c, pkg.PermissionModeration) {
	// 	return
	// }

	var botChannels []pkg.BotChannel
	var channelName string

	for it := c.Application.TwitchBots().Iterate(); it.Next(); {
		bot := it.Value()
		if bot == nil {
			fmt.Println("nil bot DansGame xD")
			utils.WebWriteError(w, 500, "nil bot or something lul")
			return
		}

		botChannel := bot.GetBotChannelByID(c.Channel.GetID())

		if botChannel != nil {
			botChannels = append(botChannels, botChannel)
			channelName = botChannel.Channel().GetName()
		}
	}

	if len(botChannels) == 0 {
		utils.WebWriteError(w, 404, "no bots in this channel?")
		return
	}

	rawMessage := "@badges=subscriber/36;color=#CC44FF;display-name=p;emotes=;flags=6-10:S.5;id=ccee20ef-9e0b-43eb-abe8-4fe9ae98e411;mod=0;room-id=11148817;subscriber=0;tmi-sent-ts=1551384658800;turbo=0;user-id=245004819;user-type= :p!p@p.tmi.twitch.tv PRIVMSG #" + channelName + " :" + message

	msg := twitch.ParseMessage(rawMessage)

	privMsg, ok := msg.(*twitch.PrivateMessage)
	if !ok {
		utils.WebWriteError(w, 404, "invalid message")
		return
	}

	response := CheckMessageSuccessResponse{
		Banned:       false,
		InputMessage: message,
	}

	for _, botChannel := range botChannels {
		event := pkg.MessageEvent{
			BaseEvent: pkg.BaseEvent{
				UserStore: botChannel.Bot().GetUserStore(),
			},
			User:    users.NewTwitchUser(privMsg.User, privMsg.Tags["user-id"]),
			Message: pb2twitch.NewTwitchMessage(privMsg),
			Channel: botChannel.Channel(),
		}

		// run message through modules on all bot channels until we detect an issue
		fmt.Println("Check message:", botChannel.ChannelName())
		actions := botChannel.OnModules(func(module pkg.Module) pkg.Actions {
			if module.Type() != pkg.ModuleTypeFilter {
				return nil
			}

			return module.OnMessage(event)
		})

		for _, action := range actions {
			if action == nil {
				log.Println("ACTION SHOULD NOT BE NIL HERE!!!!!!!!!!!!!!")
				continue
			}

			for _, mute := range action.Mutes() {
				response.Banned = true
				response.FilterData = append(response.FilterData, filterData{
					MuteType: mute.Type(),
					Reason:   mute.Reason(),
				})
			}
		}
	}

	utils.WebWrite(w, response)
}
