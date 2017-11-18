package web

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"

	twitch "github.com/gempir/go-twitch-irc"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pajlada/pajbot2/common/config"
	"github.com/pajlada/pajbot2/redismanager"
	"github.com/pajlada/pajbot2/sqlmanager"
)

// Config xD
type Config struct {
	Redis *redismanager.RedisManager
	SQL   *sqlmanager.SQLManager
	Bots  map[string]*twitch.Client
}

// Boss xD
type Boss struct {
	Host   string
	WSHost string
}

var (
	bots  map[string]*twitch.Client
	redis *redismanager.RedisManager
	sql   *sqlmanager.SQLManager
	hooks map[string]struct {
		Secret string `json:"secret"`
	}
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Init returns a webBoss which hosts the website
func Init(config *config.Config, webCfg *Config) *Boss {
	twitchBotOauthConfig.RedirectURL = config.Auth.Twitch.Bot.RedirectURI
	twitchBotOauthConfig.ClientID = config.Auth.Twitch.Bot.ClientID
	twitchBotOauthConfig.ClientSecret = config.Auth.Twitch.Bot.ClientSecret
	twitchBotOauthConfig.Scopes = []string{
		"user_read",
		"chat_login",
	}
	twitchBotOauthConfig.Endpoint = oauth2.Endpoint{
		AuthURL:  "https://api.twitch.tv/kraken/oauth2/authorize",
		TokenURL: "https://api.twitch.tv/kraken/oauth2/token",
	}
	twitchUserOauthConfig.RedirectURL = config.Auth.Twitch.User.RedirectURI
	twitchUserOauthConfig.ClientID = config.Auth.Twitch.User.ClientID
	twitchUserOauthConfig.ClientSecret = config.Auth.Twitch.User.ClientSecret
	twitchUserOauthConfig.Scopes = []string{
		"user_read",
		"channel_commercial",
		"channel_subscriptions",
		"channel_check_subscription",
		"channel_feed_read",
		"channel_feed_edit",
	}
	twitchUserOauthConfig.Endpoint = oauth2.Endpoint{
		AuthURL:  "https://api.twitch.tv/kraken/oauth2/authorize",
		TokenURL: "https://api.twitch.tv/kraken/oauth2/token",
	}
	b := &Boss{
		Host:   config.WebHost,
		WSHost: "ws://" + config.WebDomain + "/ws",
	}
	bots = webCfg.Bots
	redis = webCfg.Redis
	sql = webCfg.SQL
	hooks = config.Hooks
	return b
}

// Run xD
func (b *Boss) Run() {
	// start the hub
	go Hub.run()

	r := mux.NewRouter()
	r.HandleFunc("/ws/{type}", b.wsHandler)
	r.HandleFunc("/", b.rootHandler)
	r.HandleFunc("/dashboard", b.dashboardHandler)
	// i would like to use a subdomain for this but it might be annoying for you pajaHop
	r.HandleFunc("/api", apiRootHandler)
	api := r.PathPrefix("/api").Subrouter()

	// Serve files statically from ./web/static in /static
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("web/static/"))))

	log.Printf("Starting web on host %s", b.Host)
	InitAPI(api)
	err := http.ListenAndServe(b.Host, r)
	if err != nil {
		log.Fatal(err)
	}
}

func (b *Boss) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "<h1>xD</h1>")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (b *Boss) wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("wsHandler...")
	vars := mux.Vars(r)
	messageTypeString := vars["type"]
	messageType := MessageTypeNone
	switch messageTypeString {
	case "clr":
		messageType = MessageTypeCLR
	case "dashboard":
		messageType = MessageTypeDashboard
	}

	if messageType == MessageTypeNone {
		http.Error(w, "Invalid url. Valid urls: /ws/clr and /ws/dashboard", http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		log.Printf("Upgrader error: %v", err)
		return
	}

	log.Println("Got message!")

	// Create a custom connection
	conn := &WSConn{
		send:        make(chan []byte, 256),
		ws:          ws,
		messageType: messageType,
	}
	Hub.register <- conn
	go conn.writePump()
	conn.readPump()
}
