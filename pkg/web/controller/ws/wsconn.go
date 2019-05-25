package ws

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/users"
	"github.com/pajbot/utils"
	"github.com/pajbot/pajbot2/pkg/web/state"
	"github.com/tevino/abool"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

// WSConn xD
type WSConn struct {
	ws         *websocket.Conn
	send       chan []byte
	connected_ *abool.AtomicBool

	messageType MessageType

	c state.State

	// user is nil if the user has not authenticated
	user pkg.User
}

var _ pkg.PubSubConnection = &WSConn{}
var _ pkg.PubSubSource = &WSConn{}

func NewWSConn(ws *websocket.Conn, messageType MessageType, c state.State) *WSConn {
	conn := &WSConn{
		send:        make(chan []byte, 256),
		connected_:  abool.New(),
		ws:          ws,
		messageType: messageType,
		c:           c,
	}

	if !conn.authenticate() {
		fmt.Println("ws conn could not authenticate")
	}

	return conn
}

type pubsubMessage struct {
	Type  string
	Topic string
	Data  json.RawMessage
}

func (c *WSConn) Connection() pkg.PubSubConnection {
	return c
}

func (c *WSConn) AuthenticatedUser() pkg.User {
	return c.user
}

func (c *WSConn) IsApplication() bool {
	return false
}

func (c *WSConn) MessageReceived(source pkg.PubSubSource, topic string, bytes []byte) error {
	if !c.connected() {
		return errors.New("Connection no longer connected")
	}

	if !source.IsApplication() {
		fmt.Println("Cannot forward messages to websocket listener from anything other than the application")
		return nil
	}

	msg := pubsubMessage{
		Type:  "Publish",
		Topic: topic,
		Data:  bytes,
	}

	msgBytes, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println("error marshalling pubsub message", err)
		return nil
	}

	select {
	case c.send <- msgBytes:
	default:
		return errors.New("Connection no longer connected")
	}

	return nil
}

func (c *WSConn) pongHandler(string) error {
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	return nil
}

func (c *WSConn) connected() bool {
	return c.connected_.IsSet()
}

func (c *WSConn) authenticate() bool {
	if c.c.Session == nil {
		return false
	}

	c.user = users.NewSimpleTwitchUser(c.c.Session.TwitchUserID, c.c.Session.TwitchUserName)

	return true
}

func (c *WSConn) disconnect() {
	c.connected_.UnSet()
	close(c.send)
}

func (c *WSConn) onConnected() {
	c.connected_.Set()
	go c.writePump()
	c.readPump()
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
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				fmt.Printf("error: %v\n", err)
			}
			fmt.Printf("Error reading message: %v\n", err)
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, utils.LineFeed, utils.Space, -1))

		fmt.Println("Got msg:", string(message))

		err = c.c.PubSub.HandleJSON(c, message)
		if err != nil {
			fmt.Println("Error calling HandleJSON for pubsub:", err)
		}
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
				w.Write(utils.CRLF)
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
