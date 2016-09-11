package web

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/apirequest"
	"github.com/pajlada/pajbot2/bot"
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
		log.Error(err)
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
	var bot *bot.Bot
	var ok bool
	var p interface{}
	for _, _bots := range bots {
		if bot, ok = _bots[channel]; !ok {
			p = apiError{
				Err: "channel not found",
			}
		} else {
			ep, _rest := getEndPoint(v["rest"])
			p = exec(channel, ep, _rest)
		}
	}
	log.Debug(p)
	write(w, p)
	//p.Write(w)
	log.Info(bot != nil)
	//bot.Say("LUL")
}

func apiHook(w http.ResponseWriter, r *http.Request) {
	p := customPayload{}
	v := mux.Vars(r)
	hookType := r.Header.Get("x-github-event")
	hookSignature := r.Header.Get("x-hub-signature")
	channel := v["channel"]

	// Get hook from config according to channel
	channelHook, ok := hooks[channel]
	if !ok {
		// No hook for this channel found
		p.Add("error", "No hook found for given channel")
		write(w, p.data)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.Add("error", "Internal error")
		write(w, p.data)
		return
	}

	verified := verifySignature(channelHook.Secret, hookSignature, body)

	if !verified {
		p.Add("error", "Invalid secret")
		write(w, p.data)
		return
	}

	switch hookType {
	case "push":
		var pushData PushHookResponse
		err = json.Unmarshal(body, &pushData)
		if err != nil {
			p.Add("error", "Json Unmarshal error: "+err.Error())
			write(w, p.data)
			return
		}
		var b *bot.Bot
		for _, botMap := range bots {
			b, ok = botMap[channel]
			if ok {
				break
			}
		}
		if b == nil {
			// no bot found for channel
			p.Add("error", "No bot found for channel "+channel)
			write(w, p.data)
			return
		}

		delay := 100

		for _, commit := range pushData.Commits {
			//time.AfterFunc(time.Millisecond*time.Duration(delay), func() { writeCommit(b, commit, pushData.Repository) })
			writeCommit(b, commit, pushData.Repository)
			delay += 100
		}
		p.Add("success", true)
	}
	for _, botList := range bots {
		for key, bot := range botList {
			log.Debug(key)
			log.Debug(bot)
		}
	}

	write(w, p.data)
}

func writeCommit(b *bot.Bot, commit Commit, repository RepositoryData) {
	msg := fmt.Sprintf("%s (%s) committed to %s (%s): %s", commit.Author.Name, commit.Author.Username, repository.Name, commit.Timestamp, commit.Message)
	log.Debug(msg)
	b.SaySafef(msg)
}

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

func verifySignature(secretString string, signature string, body []byte) bool {
	const signaturePrefix = "sha1="
	const signatureLength = 45

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	secret := []byte(secretString)
	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	return hmac.Equal(signBody(secret, body), actual)
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
		log.Errorf("Code exchange failed with %s", err)
	}

	requestParameters := url.Values{}

	p := customPayload{}

	var data twitchKrakenOauth

	onSuccess := func() {
		p.Add("data", data)

		if data.Identified && data.Token.Valid {
			p.Add("username", data.Token.UserName)
			p.Add("token", token.AccessToken)
			p.Add("refreshtoken", token.RefreshToken)
			common.CreateBotAccount(sql.Session, data.Token.UserName, token.AccessToken, token.RefreshToken)
		}
	}

	apirequest.Twitch.Get("/", requestParameters, token.AccessToken, &data, onSuccess, onHTTPError, onInternalError)

	// We should, instead of returning the data raw, do something about it.
	// Right now this is useful for new apps that need access.
	// oo, do we keep multiple applications? One for bot accounts, one for clients? yes I think that sounds good
	write(w, p.data)

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
		log.Errorf("Code exchange failed with %s", err)
	}

	requestParameters := url.Values{}

	var data twitchKrakenOauth

	onSuccess := func() {
		p.Add("data", data)

		if data.Identified && data.Token.Valid {
			p.Add("username", data.Token.UserName)
			p.Add("token", token.AccessToken)
			p.Add("refreshtoken", token.RefreshToken)
			common.CreateDBUser(sql.Session, data.Token.UserName, token.AccessToken, token.RefreshToken)
		}
	}

	apirequest.Twitch.Get("/", requestParameters, token.AccessToken, &data, onSuccess, onHTTPError, onInternalError)

	// We should, instead of returning the data raw, do something about it.
	// Right now this is useful for new apps that need access.
	// oo, do we keep multiple applications? One for bot accounts, one for clients? yes I think that sounds good
	write(w, p.data)
}

func onHTTPError(statusCode int, statusMessage, errorMessage string) {
	log.Debug("HTTPERROR")
}

func onInternalError(err error) {
	log.Debugf("internal error: %s", err)
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

type twitchKrakenOauth struct {
	Identified bool `json:"identified"`
	Links      struct {
		User     string `json:"user"`
		Channel  string `json:"channel"`
		Search   string `json:"search"`
		Streams  string `json:"streams"`
		Ingests  string `json:"ingests"`
		Teams    string `json:"teams"`
		Users    string `json:"users"`
		Channels string `json:"channels"`
		Chat     string `json:"chat"`
	} `json:"_links"`
	Token struct {
		Valid         bool `json:"valid"`
		Authorization struct {
			Scopes    []string  `json:"scopes"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"authorization"`
		UserName string `json:"user_name"`
		ClientID string `json:"client_id"`
	} `json:"token"`
}
