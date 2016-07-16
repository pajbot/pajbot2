package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pajlada/pajbot2/common/config"
)

// Boss xD
type Boss struct {
	Host   string
	WSHost string
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Init returns a webBoss which hosts the website
func Init(config *config.Config) *Boss {
	b := &Boss{
		Host:   config.WebHost,
		WSHost: "ws://" + config.WebDomain + "/ws",
	}
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

	// Serve files statically from ./web/static in /static
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("web/static/"))))

	log.Infof("Starting web on host %s", b.Host)
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
		log.Errorf("Upgrader error: %v", err)
		return
	}

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
