package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pajlada/pajbot2/common"
)

// Boss xD
type Boss struct {
	Host string
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Init returns a webBoss which hosts the website
func Init(config *common.Config) *Boss {
	webHost := ":2355" // TEMPORARY
	boss := &Boss{
		Host: webHost,
	}
	return boss
}

// Run xD
func (boss *Boss) Run() {
	// start the hub
	go hub.run()

	/*
		r := mux.NewRouter()
		r.HandleFunc("/ws", boss.wsHandler)
		r.HandleFunc("/", boss.rootHandler)
		r.HandleFunc("/dashboard", boss.dashboardHandler)
	*/
	http.HandleFunc("/ws", boss.wsHandler)
	http.HandleFunc("/", boss.rootHandler)
	http.HandleFunc("/dashboard", boss.dashboardHandler)

	log.Infof("Starting web on host %s", boss.Host)
	err := http.ListenAndServe(boss.Host, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (boss *Boss) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "<h1>xD</h1>")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (boss *Boss) wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		log.Errorf("Upgrader error: %v", err)
		return
	}

	// Create a custom connection
	conn := &WSConn{send: make(chan []byte, 256), ws: ws}
	hub.register <- conn
	go conn.writePump()
	conn.readPump()
}
