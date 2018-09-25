package web

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/common/config"
	"github.com/pajlada/pajbot2/pkg/pubsub"
)

// Config xD
type Config struct {
	Redis *redis.Pool
	SQL   *sql.DB
	Bots  map[string]pkg.Sender
}

// Boss xD
type Boss struct {
	Host   string
	WSHost string
}

var (
	router      *mux.Router
	twitchBots  map[string]pkg.Sender
	redisClient *redis.Pool
	sqlClient   *sql.DB
	hooks       map[string]struct {
		Secret string
	}
	pubSub          *pubsub.PubSub
	twitchUserStore pkg.UserStore
)

var (
	newline = []byte{'\n'}
	crlf    = []byte("\r\n")
	space   = []byte{' '}
)

// Init returns a webBoss which hosts the website
func Init(config *config.Config, webCfg *Config, _pubSub *pubsub.PubSub, _twitchUserStore pkg.UserStore) *Boss {
	pubSub = _pubSub
	twitchUserStore = _twitchUserStore
	twitchBotOauth.RedirectURL = config.Auth.Twitch.Bot.RedirectURI
	twitchBotOauth.ClientID = config.Auth.Twitch.Bot.ClientID
	twitchBotOauth.ClientSecret = config.Auth.Twitch.Bot.ClientSecret
	twitchBotOauth.Scopes = []string{
		"user_read",
		"chat_login",
	}
	twitchBotOauth.Endpoint = oauth2.Endpoint{
		AuthURL:  "https://id.twitch.tv/oauth2/authorize",
		TokenURL: "https://id.twitch.tv/oauth2/token",
	}
	twitchStreamerOauth.RedirectURL = config.Auth.Twitch.Streamer.RedirectURI
	twitchStreamerOauth.ClientID = config.Auth.Twitch.Streamer.ClientID
	twitchStreamerOauth.ClientSecret = config.Auth.Twitch.Streamer.ClientSecret
	twitchStreamerOauth.Scopes = []string{
		"user_read",
		"channel_commercial",
		"channel_subscriptions",
		"channel_check_subscription",
		"channel_feed_read",
		"channel_feed_edit",
	}
	twitchStreamerOauth.Endpoint = oauth2.Endpoint{
		AuthURL:  "https://id.twitch.tv/oauth2/authorize",
		TokenURL: "https://id.twitch.tv/oauth2/token",
	}

	twitchUserOauth.RedirectURL = config.Auth.Twitch.User.RedirectURI
	twitchUserOauth.ClientID = config.Auth.Twitch.User.ClientID
	twitchUserOauth.ClientSecret = config.Auth.Twitch.User.ClientSecret
	twitchUserOauth.Scopes = []string{
		"openid",
	}
	twitchUserOauth.Endpoint = oauth2.Endpoint{
		AuthURL:  "https://id.twitch.tv/oauth2/authorize",
		TokenURL: "https://api.twitch.tv/kraken/oauth2/token",
	}
	b := &Boss{
		Host:   config.Web.Host,
		WSHost: "ws://" + config.Web.Domain + "/ws",
	}
	twitchBots = webCfg.Bots
	redisClient = webCfg.Redis
	sqlClient = webCfg.SQL
	hooks = config.Hooks

	router = mux.NewRouter()
	router.HandleFunc("/", b.rootHandler)
	router.HandleFunc("/ws/{type}", b.wsHandler)
	router.HandleFunc("/dashboard", b.dashboardHandler)
	// i would like to use a subdomain for this but it might be annoying for you pajaHop
	router.HandleFunc("/api", apiRootHandler)
	api := router.PathPrefix("/api").Subrouter()

	initAPI(api, config)

	// Serve files statically from ./web/static in /static
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("web/static/"))))

	fmt.Printf("Starting web on host %s\n", b.Host)
	return b
}

// Run xD
func (b *Boss) Run() {
	go Hub.run()

	corsObj := handlers.AllowedOrigins([]string{"*"})
	err := http.ListenAndServe(b.Host, handlers.CORS(corsObj)(router))
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
	fmt.Println("ws handler")
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
		fmt.Println("ws handler error")
		http.Error(w, "Invalid url. Valid urls: /ws/clr and /ws/dashboard", http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		fmt.Printf("Upgrader error: %v\n", err)
		return
	}

	conn := NewWSConn(ws, messageType)

	// Create a custom connection
	Hub.register <- conn
	conn.onConnected()
}
