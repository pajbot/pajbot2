package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"time"

	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/pkg/common/config"
)

// api endpoints
const (
	// single user
	USER = "user"
	// list of users, might be useful for rank lists or stuff like that
	USERS = "users"
	// single command by id or trigger
	COMMAND = "command"
	// list of all commands
	COMMANDS = "commands"
	// single module by ID
	MODULE = "module"
	// list of modules
	MODULES = "modules"
	// single banphrase by ID
	BANPHRASE = "banphrase"
	// list of all banphrases
	BANPHRASES = "banphrases"
)

func newError(err string) interface{} {
	return apiError{
		Err: err,
	}
}

func writeError(w http.ResponseWriter, message string) {
	data := struct {
		Message string
	}{
		message,
	}
	bs, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error in web write: %s", err)
		bs, _ = json.Marshal(newError("internal server error"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

func write(w http.ResponseWriter, data interface{}) {
	bs, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error in web write: %s", err)
		bs, _ = json.Marshal(newError("internal server error"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

func users(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	w.Write([]byte(v["user"]))
}

func getEndPoint(url string) (string, string) {
	spl := strings.SplitN(url, "/", 2)
	var spl1 string
	if len(spl) > 1 {
		spl1 = spl[1]
	}
	return strings.ToLower(spl[0]), spl1
}

func exec(channel, endpoint, rest string) interface{} {
	if !isValidURL(rest) {
		return newError(ErrInvalidUserName)
	}
	var p interface{}
	switch endpoint {
	case USER:
		if !isValidUserName(rest) {
			return newError(ErrInvalidUserName)
		}
		p = newError("unimplemented")
	default:
		p = newError("invalid endpoint")
	}
	return p
}

// APIHandler xD
func APIHandler(w http.ResponseWriter, r *http.Request) {
	return
	/*
		// TODO: check user auth and permissions
		// TODO: route http methods
		v := mux.Vars(r)
		channel := v["channel"]
		var bot *bots.TwitchBot
		var ok bool
		var p interface{}
		if bot, ok = twitchBots[channel]; !ok {
			p = apiError{
				Err: "channel not found",
			}
		} else {
			ep, _rest := getEndPoint(v["rest"])
			p = exec(channel, ep, _rest)
		}
		fmt.Printf("Bot: %#v", bot)
		write(w, p)
	*/
}

func apiRootHandler(w http.ResponseWriter, r *http.Request) {
	p := customPayload{}
	p.Add("paja", "Dank")
	write(w, p.data)
}

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
	UserID    string
	Action    string
	Duration  int
	TargetID  string
	Reason    string
	Timestamp time.Time
	Context   *string
}

type moderationResponse struct {
	ChannelID string

	Actions []moderationAction
}

func apiChannelModerationLatest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response := moderationResponse{}

	response.ChannelID = vars["channelID"]

	const queryF = "SELECT `UserID`, `Action`, `Duration`, `TargetID`, `Reason`, `Timestamp`, `Context` FROM `ModerationAction` WHERE `ChannelID`=? ORDER BY `Timestamp` DESC LIMIT 20;"

	rows, err := sqlClient.Query(queryF, response.ChannelID)
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

	write(w, response)
}

func onHTTPError(statusCode int, statusMessage, errorMessage string) {
	log.Println("HTTPERROR: ", errorMessage)
}

func onInternalError(err error) {
	fmt.Printf("internal error: %s", err)
}

func initAPI(m *mux.Router, config *config.Config) {
	m.HandleFunc("/", apiRootHandler)

	fmt.Println("Init API")
	err := twitchAuthInit(m, &config.Auth.Twitch)
	if err != nil {
		panic(err)
	}

	// m.HandleFunc(`/channel/{channel:\w+}/{rest:.*}`, APIHandler)
	m.HandleFunc(`/channel/{channelID}/moderation/latest`, apiChannelModerationLatest)
	m.HandleFunc(`/hook/{channel:\w+}`, apiHook)
	m.HandleFunc(`/callbacks/follow`, apiCallbacksFollow)
	m.HandleFunc(`/callbacks/streams`, apiCallbacksStreams)
}
