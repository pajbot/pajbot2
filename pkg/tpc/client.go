package tpc

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/scylladb/go-set/u64set"
)

type TweetProviderClient struct {
	writer        chan []byte
	closer        chan bool
	readerDone    chan struct{}
	subscriptions *u64set.Set
	conn          *websocket.Conn
	onAck         func(...uint64)
}

func New() *TweetProviderClient {
	return &TweetProviderClient{
		writer:        make(chan []byte, 100),
		closer:        make(chan bool),
		readerDone:    make(chan struct{}),
		subscriptions: u64set.New(),
	}
}

func (c *TweetProviderClient) OnAck(cb func(userIDs ...uint64)) {
	c.onAck = cb
}

func (c *TweetProviderClient) readPump() {
	c.readerDone = make(chan struct{})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			close(c.readerDone)
			return
		}

		// parsemessage
		var msg incomingMessage
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Println("Error unmarshalling:", err)
			return
		}
		if msg.Type == "ack_subscriptions" {
			var userIDs []uint64
			json.Unmarshal(msg.Data, &userIDs)
			if c.onAck != nil {
				c.onAck(userIDs...)
			}
		}
	}

}

func (c *TweetProviderClient) writePump() {
	ping := time.NewTicker(time.Second)
	defer ping.Stop()

	for {
		select {
		case msg := <-c.writer:
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-ping.C:
			err := c.conn.WriteMessage(websocket.PingMessage, []byte("pb2"))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-c.closer:
			err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Closing"))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-c.readerDone:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (c *TweetProviderClient) Start() error {
	if c.conn != nil {
		return ErrAlreadyConnected
	}
	var err error
	c.conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:2356", nil)
	if err != nil {
		fmt.Println("Error in dial:", err)
		return err
	}

	go c.readPump()
	go c.writePump()

	return nil
}

func (c *TweetProviderClient) Close() error {
	if c.conn == nil {
		return ErrAlreadyDisconnected
	}
	c.closer <- true
	<-c.readerDone
	c.conn.Close()
	c.conn = nil

	return nil
}

// InsertSubscriptions xd
func (c *TweetProviderClient) InsertSubscriptions(userIDs ...uint64) error {
	c.subscriptions.Add(userIDs...)

	msg := outgoingMessage{
		Type: "insert_subscriptions",
		Data: userIDs,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.writer <- bytes

	return nil
}

// RemoveSubscriptions xd
func (c *TweetProviderClient) RemoveSubscriptions(userIDs []uint64) error {
	c.subscriptions.Remove(userIDs...)

	msg := outgoingMessage{
		Type: "remove_subscriptions",
		Data: userIDs,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.writer <- bytes

	return nil
}
