package moderation

import (
	"fmt"
	"net/http"

	twitch "github.com/gempir/go-twitch-irc/v2"
	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/pkg"
	pb2twitch "github.com/pajlada/pajbot2/pkg/twitch"
	"github.com/pajlada/pajbot2/pkg/users"
	"github.com/pajlada/pajbot2/pkg/utils"
	"github.com/pajlada/pajbot2/pkg/web/state"
	_ "github.com/swaggo/echo-swagger" // echo-swagger middleware
)

func apiCheckMessageMissingVariables(w http.ResponseWriter, r *http.Request) {
	utils.WebWriteError(w, 400, "Missing message parameter")
}

type filterData struct {
	ActionType pkg.ActionType `json:"action_type"`
	Reason     string         `json:"reason"`
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

	channelID := vars["channelID"]
	message := vars["message"]

	if channelID == "" {
		// ERROR: Missing required ChannelID parameter
		utils.WebWriteError(w, 400, "Missing required channelID parameter")
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

	fmt.Println("Channel ID:", channelID)

	for it := c.Application.TwitchBots().Iterate(); it.Next(); {
		bot := it.Value()
		if bot == nil {
			fmt.Println("nil bot DansGame xD")
			utils.WebWriteError(w, 500, "nil bot or something lul")
			return
		}

		botChannel := bot.GetBotChannelByID(channelID)

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
		// run message through modules on all bot channels until we detect an issue
		fmt.Println("Check message:", botChannel.ChannelName())
		lolMessage := pb2twitch.NewTwitchMessage(privMsg)
		botChannel.OnModules(func(module pkg.Module) error {
			if module.Spec().Type() != pkg.ModuleTypeFilter {
				return nil
			}
			action := &pkg.TwitchAction{
				Sender:  botChannel.Bot(),
				Channel: botChannel.Channel(),
				User:    users.NewTwitchUser(privMsg.User, privMsg.Tags["user-id"]),
				Soft:    true,
			}
			module.OnMessage(botChannel, action.User, lolMessage, action)
			if action.Get() != nil {
				response.Banned = true
				response.FilterData = append(response.FilterData, filterData{
					ActionType: action.Get().Type(),
					Reason:     action.Reason(),
				})
				// return errors.New("stop here xd")
			}
			return nil
		})
	}

	utils.WebWrite(w, response)
}
