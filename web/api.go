package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/bot"
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

func write(w http.ResponseWriter, data interface{}) {
	bs, err := json.Marshal(data)
	if err != nil {
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
	log.Info(channel, endpoint, rest)
	var p interface{}
	switch endpoint {
	case USER:
		p = getUserPayload(channel, rest)
	default:
		p = newError("invalid endpoint")
	}
	return p
}

// APIHandler xD
func APIHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: check user auth and permissions
	// TODO: route http methods
	v := mux.Vars(r)
	channel := v["channel"]
	var bot *bot.Bot
	var ok bool
	var p interface{}
	if bot, ok = bots[channel]; !ok {
		p = apiError{
			Err: "channel not found",
		}
	} else {
		ep, _rest := getEndPoint(v["rest"])
		p = exec(channel, ep, _rest)
	}
	log.Debug(p)
	write(w, p)
	//p.Write(w)
	log.Info(bot != nil)
	//bot.Say("LUL")
}

func apiRootHandler(w http.ResponseWriter, r *http.Request) {
	p := customPayload{}
	p.Add("paja", "Dank")
	write(w, p.data)
}

// InitAPI adds routes to the given subrouter
func InitAPI(m *mux.Router) {
	m.HandleFunc("/", apiRootHandler)
	m.HandleFunc(`/{channel:\w+}/{rest:.*}`, APIHandler)
}
