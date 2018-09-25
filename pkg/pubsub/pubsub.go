package pubsub

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/pajlada/pajbot2/pkg"
)

type topic struct {
}

type Connection interface {
	MessageReceived(topic string, bytes []byte, Authorization *pkg.PubSubAuthorization) error
}

type SubscriptionHandler interface {
	ConnectionSubscribed(connection Connection, topic string, Authorization *pkg.PubSubAuthorization) (error, bool)
}

type SubscriptionType int

const (
	SubscriptionTypeOnce SubscriptionType = iota
	SubscriptionTypeContinuous
)

const SubscribeOnce = 0

type operationType int

const (
	operationPublish operationType = iota
	operationSubscribe
)

type Message struct {
	operation     operationType
	topic         string
	data          interface{}
	authorization *pkg.PubSubAuthorization
	connection    Connection
}

type PubSub struct {
	c           chan (Message)
	connections []Connection

	topicsMutex sync.Mutex
	topics      map[string][]*Listener

	onSubscribeMutex sync.Mutex
	onSubscribe      map[string][]SubscriptionHandler
}

func New() *PubSub {
	return &PubSub{
		c:           make(chan (Message)),
		topics:      make(map[string][]*Listener),
		onSubscribe: make(map[string][]SubscriptionHandler),
	}
}

func (ps *PubSub) AcceptConnection(conn Connection) {
	ps.connections = append(ps.connections, conn)
}

func (ps *PubSub) Run() {
	for {
		select {
		case msg := <-ps.c:
			switch msg.operation {
			case operationPublish:
				ps.publish(msg.topic, msg.data, msg.authorization)
			case operationSubscribe:
				ps.Subscribe(msg.connection, msg.topic, msg.authorization)

			default:
				fmt.Printf("Unhandled operation: %v\n", msg.operation)
			}
		}
	}
}

func (ps *PubSub) publish(topic string, data interface{}, auth *pkg.PubSubAuthorization) {
	ps.topicsMutex.Lock()

	for _, l := range ps.topics[topic] {
		go func(listener *Listener) {
			err := listener.Publish(topic, data, auth)
			if err != nil {
				ps.UnsubscribeAll(listener)
				fmt.Println(err)
			}
		}(l)
	}

	ps.topicsMutex.Unlock()
}

func (ps *PubSub) Publish(topic string, data interface{}, auth *pkg.PubSubAuthorization) {
	ps.c <- Message{
		operation:     operationPublish,
		topic:         topic,
		data:          data,
		authorization: auth,
	}
}

type pubsubMessage struct {
	Type          string
	Topic         string
	Data          interface{} `json:",omitempty"`
	Authorization *pkg.PubSubAuthorization
}

func (ps *PubSub) HandleJSON(connection Connection, bytes []byte) error {
	var msg pubsubMessage
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		return err
	}

	switch msg.Type {
	case "Publish":
		fmt.Println("Send publish message on topic", msg.Topic)
		ps.c <- Message{operation: operationPublish, topic: msg.Topic, connection: connection, data: msg.Data, authorization: msg.Authorization}
	case "Subscribe":
		ps.c <- Message{operation: operationSubscribe, topic: msg.Topic, connection: connection, authorization: msg.Authorization}
	}

	return nil
}

func (ps *PubSub) Subscribe(connection Connection, topic string, authorization *pkg.PubSubAuthorization) {
	successfulAuthorization := ps.notifySubscriptionHandlers(connection, topic, authorization)
	if !successfulAuthorization {
		fmt.Println("Failed authorization")
		return
	}

	{
		ps.topicsMutex.Lock()
		defer ps.topicsMutex.Unlock()
		l := &Listener{connection, SubscriptionTypeContinuous}

		ps.topics[topic] = append(ps.topics[topic], l)
	}
}

func (ps *PubSub) notifySubscriptionHandlers(connection Connection, topic string, authorization *pkg.PubSubAuthorization) bool {
	ps.onSubscribeMutex.Lock()
	defer ps.onSubscribeMutex.Unlock()

	for _, handler := range ps.onSubscribe[topic] {
		err, successfulAuthorization := handler.ConnectionSubscribed(connection, topic, authorization)
		if err != nil {
			fmt.Println("Error in subscription handler:", err)
		}

		if !successfulAuthorization {
			return false
		}
	}

	return true
}

func (ps *PubSub) HandleSubscribe(connection SubscriptionHandler, topic string) {
	ps.onSubscribeMutex.Lock()
	defer ps.onSubscribeMutex.Unlock()

	ps.onSubscribe[topic] = append(ps.onSubscribe[topic], connection)
}

func (ps *PubSub) SubscribeOnce(connection Connection, topic string) {
	ps.topicsMutex.Lock()
	defer ps.topicsMutex.Unlock()

	ps.topics[topic] = append(ps.topics[topic], &Listener{connection, SubscriptionTypeOnce})
}

func (ps *PubSub) UnsubscribeAll(l *Listener) {
	ps.topicsMutex.Lock()
	defer ps.topicsMutex.Unlock()

	for topic, listeners := range ps.topics {
		var newListeners []*Listener
		for _, listener := range listeners {
			if listener != l {
				newListeners = append(newListeners, listener)
			}
		}

		ps.topics[topic] = newListeners
	}
}
