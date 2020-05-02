package moderation

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/users"
	"github.com/pajbot/pajbot2/pkg/web/state"
	"github.com/pajbot/utils"
)

const ActionUnknown = 0
const ActionTimeout = 1
const ActionBan = 2
const ActionUnban = 3

func getActionString(action int) string {
	switch action {
	case ActionTimeout:
		return "timeout"

	case ActionBan:
		return "ban"

	case ActionUnban:
		return "unban"
	}

	return ""
}

type moderationAction struct {
	UserName   string
	UserID     string
	Action     string
	Duration   int
	TargetID   string
	TargetName string
	Reason     string
	Timestamp  time.Time
	Context    *string
}

type moderationResponse struct {
	ChannelID string

	Actions []moderationAction
}

func apiChannelModerationLatest(w http.ResponseWriter, r *http.Request) {
	c := state.Context(w, r)

	if c.Channel == nil {
		utils.WebWriteError(w, 500, "this is not a channel we are in")
		return
	}

	if c.Session == nil {
		utils.WebWriteError(w, 400, "Not authorized to view this endpoint")
		return
	}

	user := users.NewSimpleTwitchUser(c.Session.TwitchUserID, c.Session.TwitchUserName)
	if user == nil {
		utils.WebWriteError(w, 400, "Not authorized to view this endpoint")
		return
	}

	if !user.HasPermission(c.Channel, pkg.PermissionModeration) && !user.HasPermission(c.Channel, pkg.PermissionReport) && !user.HasPermission(c.Channel, pkg.PermissionAdmin) && !user.HasPermission(c.Channel, pkg.PermissionReportAPI) {
		utils.WebWriteError(w, 400, "Not authorized to view this endpoint!!!")
		return
	}

	vars := mux.Vars(r)
	response := moderationResponse{}

	response.ChannelID = c.Channel.GetID()

	fmt.Println("Channel ID:", vars)

	const queryF = "SELECT user_id, action, duration, target_id, reason, timestamp, context FROM moderation_action WHERE channel_id=$1 ORDER BY timestamp DESC LIMIT 20;"

	rows, err := c.SQL.Query(queryF, response.ChannelID) // GOOD
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		action := moderationAction{}
		actionIndex := 0
		if err := rows.Scan(&action.UserID, &actionIndex, &action.Duration, &action.TargetID, &action.Reason, &action.Timestamp, &action.Context); err != nil {
			panic(err)
		}
		action.Action = getActionString(actionIndex)

		response.Actions = append(response.Actions, action)
	}

	utils.WebWrite(w, response)
}
