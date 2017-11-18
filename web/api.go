package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dankeroni/gotwitch"
	twitch "github.com/gempir/go-twitch-irc"
	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/apirequest"
	"github.com/pajlada/pajbot2/common"
	"golang.org/x/oauth2"
)

var (
	twitchBotOauthConfig  = &oauth2.Config{}
	twitchUserOauthConfig = &oauth2.Config{}
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
		log.Printf("Error in web write: %s", err)
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
	var bot *twitch.Client
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
	log.Printf("Bot: %#v", bot)
	write(w, p)
}

func apiRootHandler(w http.ResponseWriter, r *http.Request) {
	p := customPayload{}
	p.Add("paja", "Dank")
	write(w, p.data)
}

// TODO(pajlada): This should be random per request
var oauthStateString = "penis"

func apiTwitchBotLogin(w http.ResponseWriter, r *http.Request) {
	url := twitchBotOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func apiTwitchBotCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("Invalid oauth state")
		// bad oauth state
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := twitchBotOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed with %s", err)
	}
	write(w, "Access token: "+token.AccessToken)

	p := customPayload{}

	onSuccess := func(data gotwitch.Self) {
		log.Println("Success!!!!!!!!!!")
		p.Add("data", data)

		if data.Identified && data.Token.Valid {
			p.Add("username", data.Token.UserName)
			p.Add("token", token.AccessToken)
			p.Add("refreshtoken", token.RefreshToken)
			err = common.CreateDBUser(sql.Session, data.Token.UserName, token.AccessToken, token.RefreshToken, "bot")
			if err != nil {
				// XXX: handle this
				log.Println(err)
			}
		}
	}

	apirequest.Twitch.GetSelf(token.AccessToken, onSuccess, onHTTPError, onInternalError)

	// We should, instead of returning the data raw, do something about it.
	// Right now this is useful for new apps that need access.
	// oo, do we keep multiple applications? One for bot accounts, one for clients? yes I think that sounds good
	write(w, p.data)
	log.Println("hehe")

	//common.CreateBotAccount(sql.Session, )
}

func apiTwitchUserLogin(w http.ResponseWriter, r *http.Request) {
	url := twitchUserOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func apiTwitchUserCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("Invalid oauth state")
		// bad oauth state
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	p := customPayload{}

	code := r.FormValue("code")
	if code == "" {
		// no valid code given
		p.Add("error", "Invalid code")
		write(w, p.data)
		return
	}
	token, err := twitchUserOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed with %s", err)
	}

	onSuccess := func(data gotwitch.Self) {
		log.Println("Success!")
		p.Add("data", data)

		if data.Identified && data.Token.Valid {
			p.Add("username", data.Token.UserName)
			p.Add("token", token.AccessToken)
			p.Add("refreshtoken", token.RefreshToken)
			err = common.CreateDBUser(sql.Session, data.Token.UserName, token.AccessToken, token.RefreshToken, "user")
			if err != nil {
				// XXX: handle this
				log.Println(err)
			}
		}
	}

	apirequest.TwitchV3.GetSelf(token.AccessToken, onSuccess, onHTTPError, onInternalError)

	// We should, instead of returning the data raw, do something about it.
	// Right now this is useful for new apps that need access.
	// oo, do we keep multiple applications? One for bot accounts, one for clients? yes I think that sounds good
	write(w, p.data)
}

func onHTTPError(statusCode int, statusMessage, errorMessage string) {
	log.Println("HTTPERROR: ", errorMessage)
}

func onInternalError(err error) {
	log.Printf("internal error: %s", err)
}

// InitAPI adds routes to the given subrouter
func InitAPI(m *mux.Router) {
	m.HandleFunc("/", apiRootHandler)
	m.HandleFunc("/auth/twitch/bot", apiTwitchBotLogin)
	m.HandleFunc("/auth/twitch/user", apiTwitchUserLogin)
	m.HandleFunc("/auth/twitch/bot/callback", apiTwitchBotCallback)
	m.HandleFunc("/auth/twitch/user/callback", apiTwitchUserCallback)
	m.HandleFunc(`/channel/{channel:\w+}/{rest:.*}`, APIHandler)
	m.HandleFunc(`/hook/{channel:\w+}`, apiHook)
}
