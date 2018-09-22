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
	twitchBots  map[string]pkg.Sender
	redisClient *redis.Pool
	sqlClient   *sql.DB
	hooks       map[string]struct {
		Secret string
	}
	pubSub *pubsub.PubSub
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Init returns a webBoss which hosts the website
func Init(config *config.Config, webCfg *Config, _pubSub *pubsub.PubSub) *Boss {
	pubSub = _pubSub
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
		Host:   config.Web.Host,
		WSHost: "ws://" + config.Web.Domain + "/ws",
	}
	twitchBots = webCfg.Bots
	redisClient = webCfg.Redis
	sqlClient = webCfg.SQL
	hooks = config.Hooks
	return b
}

// Run xD
func (b *Boss) Run() {
	go Hub.run()

	r := mux.NewRouter()
	r.HandleFunc("/", b.rootHandler)
	r.HandleFunc("/ws/{type}", b.wsHandler)
	r.HandleFunc("/dashboard", b.dashboardHandler)
	// i would like to use a subdomain for this but it might be annoying for you pajaHop
	r.HandleFunc("/api", apiRootHandler)
	api := r.PathPrefix("/api").Subrouter()

	// Serve files statically from ./web/static in /static
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("../../web/static/"))))

	fmt.Printf("Starting web on host %s\n", b.Host)
	InitAPI(api)
	corsObj := handlers.AllowedOrigins([]string{"*"})
	err := http.ListenAndServe(b.Host, handlers.CORS(corsObj)(r))
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

	// Create a custom connection
	conn := &WSConn{
		send:        make(chan []byte, 256),
		ws:          ws,
		messageType: messageType,
	}
	fmt.Println("xd")
	Hub.register <- conn
	fmt.Println("loooooooooooool")
	go conn.writePump()
	conn.readPump()
}
