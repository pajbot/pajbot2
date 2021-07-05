package moderation

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/web/state"
	"github.com/pajbot/pajbot2/pkg/webutils"
	"github.com/pajbot/utils"
)

type userResponse struct {
	ChannelID string
	TargetID  string

	Actions []*moderationAction
}

func apiUserMissingVariables(w http.ResponseWriter, r *http.Request) {
	utils.WebWriteError(w, 400, "Missing either user_id or user_name query parameter")
}

func apiUser(w http.ResponseWriter, r *http.Request) {
	const queryF = "SELECT user_id, action, duration, target_id, reason, timestamp, context FROM moderation_action WHERE channel_id=$1 AND target_id=$2 ORDER BY timestamp DESC LIMIT 20;"

	c := state.Context(w, r)
	vars := mux.Vars(r)

	if c.Channel == nil {
		utils.WebWriteError(w, 400, "Invalid channel_id specified")
		return
	}

	response := userResponse{
		Actions:   make([]*moderationAction, 0),
		ChannelID: c.Channel.GetID(),
	}

	if userID, ok := vars["user_id"]; ok {
		response.TargetID = userID
	} else if userName, ok := vars["user_name"]; ok {
		response.TargetID = c.TwitchUserStore.GetID(userName)
		if response.TargetID == "" {
			utils.WebWriteError(w, 400, "Invalid user_name")
			return
		}
	} else {
		utils.WebWriteError(w, 400, "Missing required user_id or user_name parameter")
		return
	}

	if !webutils.RequirePermission(w, c, c.Channel, pkg.PermissionModeration) {
		return
	}

	rows, err := c.SQL.Query(queryF, response.ChannelID, response.TargetID) // GOOD
	if err != nil {
		fmt.Println("error in mysql query apiUser:", err)
		utils.WebWriteError(w, 500, "Internal error")
		return
	}

	fmt.Println("Query", response.ChannelID, response.TargetID)

	defer rows.Close()

	request := pkg.NewUserStoreRequest()

	for rows.Next() {
		action := moderationAction{}
		actionIndex := 0
		if err := rows.Scan(&action.UserID, &actionIndex, &action.Duration, &action.TargetID, &action.Reason, &action.Timestamp, &action.Context); err != nil {
			fmt.Println("error when scanning row:", err)
			utils.WebWriteError(w, 500, "Internal error")
			return
		}
		action.Action = getActionString(actionIndex)

		request.AddID(action.UserID)
		request.AddID(action.TargetID)

		response.Actions = append(response.Actions, &action)
	}

	names, _ := request.Execute(c.TwitchUserStore)

	for i, action := range response.Actions {
		for id, name := range names {
			if action.UserID == id {
				response.Actions[i].UserName = name
			}
			if action.TargetID == id {
				response.Actions[i].TargetName = name
			}
		}
	}

	utils.WebWrite(w, response)
}
