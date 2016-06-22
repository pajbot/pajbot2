package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// WSHandler is a method which handles each connection
type WSHandler func(conn *websocket.Conn)

// Boss xD
type Boss struct {
	Host    string
	Handler WSHandler
}

// Init xD
func Init(host string) *Boss {
	return &Boss{
		Host: host,
	}
}

// Run starts the websocket server
func (ws *Boss) Run() {
	// TODO: figure out https. https should be a requirement
	http.HandleFunc("/ws", ws.wsHandler)
	http.HandleFunc("/", ws.rootHandler)

	http.ListenAndServe(ws.Host, nil)
}

func (ws *Boss) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "<h1>xD</h1>")
}

func (ws *Boss) wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	go ws.Handler(conn)
}

type msg struct {
	Num int
}
