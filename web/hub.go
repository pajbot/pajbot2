package web

// Hub xD
type Hub struct {
	connections map[*WSConn]bool
	broadcast   chan []byte
	unregister  chan *WSConn
	register    chan *WSConn
}

var hub = Hub{
	connections: make(map[*WSConn]bool),
	broadcast:   make(chan []byte),
	unregister:  make(chan *WSConn),
	register:    make(chan *WSConn),
}

func (h *Hub) run() {
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
					delete(hub.connections, conn)
				}
			}
		}
	}
}
