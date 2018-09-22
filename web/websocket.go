package web

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/pubsub"
)

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

	messageType MessageType

	// user is nil if the user has not authenticated
	user pkg.User
}

var _ pubsub.Connection = &WSConn{}

func (c *WSConn) MessageReceived(topic string, bytes []byte) error {
	switch topic {
	case "MessageReceived":
		c.send <- bytes
	}
	fmt.Printf("ws conn received message on topic %s: %s\n", topic, string(bytes))
	return nil
}

func (c *WSConn) pongHandler(string) error {
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	return nil
}

// TODO: Fix proper authentication
// TODO: load user from db/redis/cache
func (c *WSConn) authenticate(username string) {
	fmt.Printf("Attempting to authenticate as %s\n", username)
}

func (c *WSConn) readPump() {
	// read shit
	defer func() {
		Hub.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(c.pongHandler)
	for {
		fmt.Println("Reading messages....")
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				fmt.Printf("error: %v\n", err)
			}
			fmt.Printf("Error reading message: %v\n", err)
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		fmt.Printf("Got the message %s\n", message)

		// pubSub.Publish("penis", "aaaaaaaaa")
		// TODO: Handle incoming messages
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

// MessageType is used to help redirect a message to the proper connections
type MessageType uint8

// All available MessageTypes
const (
	MessageTypeAll MessageType = iota
	MessageTypeNone
	MessageTypeCLR
	MessageTypeDashboard
)

// WSMessage xD
type WSMessage struct {
	Channel string

	MessageType MessageType

	// LevelRequired <=0 means the message does not require authentication, otherwise
	// authentication is required and the users level must be equal to or above
	// the LevelRequired value
	LevelRequired int

	Payload *Payload
}
