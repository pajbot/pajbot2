package moderation

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/pkg/utils"
	"github.com/pajlada/pajbot2/pkg/web/state"
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

	vars := mux.Vars(r)
	response := moderationResponse{}

	response.ChannelID = vars["channelID"]

	fmt.Println("Channel ID:", vars)

	const queryF = "SELECT `UserID`, `Action`, `Duration`, `TargetID`, `Reason`, `Timestamp`, `Context` FROM `ModerationAction` WHERE `ChannelID`=? ORDER BY `Timestamp` DESC LIMIT 20;"

	rows, err := c.SQL.Query(queryF, response.ChannelID)
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
