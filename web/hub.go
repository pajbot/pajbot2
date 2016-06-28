package web

import "encoding/json"

// ConnectionHub xD
type ConnectionHub struct {
	connections map[*WSConn]bool
	broadcast   chan []byte
	unregister  chan *WSConn
	register    chan *WSConn
}

// Hub xD
var Hub = ConnectionHub{
	connections: make(map[*WSConn]bool),
	broadcast:   make(chan []byte),
	unregister:  make(chan *WSConn),
	register:    make(chan *WSConn),
}

// Payload xD
type Payload struct {
	Event string `json:"event"`
}

// ToJSON creates a json string from the payload
func (p *Payload) ToJSON() (ret []byte) {
	ret, err := json.Marshal(p)
	if err != nil {
		log.Error("Erro marshalling payload:", err)
	}
	return
}

func (h *ConnectionHub) run() {
	for {
		select {
		case conn := <-h.register:
			log.Debugf("REGISTERING %#v", conn)
			h.connections[conn] = true
		case conn := <-h.unregister:
			if _, ok := h.connections[conn]; ok {
				delete(h.connections, conn)
				close(conn.send)
			}
		case message := <-h.broadcast:
			for conn := range h.connections {
				select {
				case conn.send <- message:
				default:
					// Not sure what this is for
					close(conn.send)
					delete(h.connections, conn)
				}
			}
		}
	}
}

// Broadcast some data to all connections
func (h *ConnectionHub) Broadcast(data []byte) {
	h.broadcast <- data
}
