package moderation

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/users"
	"github.com/pajlada/pajbot2/pkg/utils"
	"github.com/pajlada/pajbot2/web/state"
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
	const queryF = "SELECT `UserID`, `Action`, `Duration`, `TargetID`, `Reason`, `Timestamp`, `Context` FROM `ModerationAction` WHERE `ChannelID`=? AND `TargetID`=? ORDER BY `Timestamp` DESC LIMIT 20;"

	c := state.Context(w, r)

	if c.Session == nil {
		utils.WebWriteError(w, 400, "Not authorized to view this endpoint")
		return
	}

	user := users.NewSimpleTwitchUser(c.Session.TwitchUserID, c.Session.TwitchUserName)
	if user == nil {
		utils.WebWriteError(w, 400, "Not authorized to view this endpoint")
		return
	}

	if !user.HasGlobalPermission(pkg.PermissionModeration) {
		utils.WebWriteError(w, 400, "Not authorized to view this endpoint!!!")
		return
	}

	vars := mux.Vars(r)
	response := userResponse{}

	response.Actions = make([]*moderationAction, 0)
	response.ChannelID = vars["channelID"]

	if response.ChannelID == "" {
		// ERROR: Missing required ChannelID parameter
		utils.WebWriteError(w, 400, "Missing required ChannelID parameter")
		return
	}

	if userID := vars["user_id"]; userID != "" {
		response.TargetID = userID
	} else if userName := vars["user_name"]; userName != "" {
		response.TargetID = c.TwitchUserStore.GetID(userName)
	}

	if response.TargetID == "" {
		// ERROR: Unable to figure out a valid user ID
		utils.WebWriteError(w, 400, "Provided user name did not return a valid user ID")
		return
	}

	rows, err := c.SQL.Query(queryF, response.ChannelID, response.TargetID)
	if err != nil {
		fmt.Println("error in mysql query apiUser:", err)
		utils.WebWriteError(w, 500, "Internal error")
		return
	}

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
