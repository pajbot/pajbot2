package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pajbot/pajbot2/pkg/web/router"
	"github.com/pajbot/pajbot2/pkg/web/state"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Load() {
	go Hub.run()

	m := router.Subrouter("/ws")

	router.RHandleFunc(m, `/{type}`, wsHandler)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
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

	c := state.Context(w, r)

	conn := NewWSConn(ws, messageType, c)

	fmt.Println("aaaaaaaaa")

	// Create a custom connection
	Hub.register <- conn
	conn.onConnected()
}
