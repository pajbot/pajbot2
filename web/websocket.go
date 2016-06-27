package web

import "github.com/gorilla/websocket"

type msg struct {
	Num int
}

// WSHandler is a method which handles each connection
type WSHandler func(conn *websocket.Conn)

func wsHandler(conn *websocket.Conn) {
	for conn != nil {
		m := msg{}

		err := conn.ReadJSON(&m)
		if err != nil {
			log.Error("Error reading json.", err)
		}

		log.Infof("Got message: %#v", m)

		if err = conn.WriteJSON(m); err != nil {
			log.Error(err)
		}
	}
}
