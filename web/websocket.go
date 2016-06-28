package web

import (
	"bytes"
	"time"

	"github.com/gorilla/websocket"
)

type msg struct {
	Num int
}

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

// WSConn xD
type WSConn struct {
	ws   *websocket.Conn
	send chan []byte
}

func (c *WSConn) pongHandler(string) error {
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	return nil
}

func (c *WSConn) readPump() {
	// read shit
	defer func() {
		hub.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(c.pongHandler)
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Errorf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Debugf("Got the message %s", message)
		hub.broadcast <- message
	}
}

func (c *WSConn) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *WSConn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued message to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
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
