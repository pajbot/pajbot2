package pubsub

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/pajbot/pajbot2/pkg"
)

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
	operation operationType
	topic     string
	data      interface{}
	source    pkg.PubSubSource
}

type PubSub struct {
	c           chan (Message)
	connections []pkg.PubSubConnection

	topicsMutex sync.Mutex
	topics      map[string][]*Listener

	onSubscribeMutex sync.Mutex
	onSubscribe      map[string][]pkg.PubSubSubscriptionHandler
}

func New() *PubSub {
	return &PubSub{
		c:           make(chan (Message)),
		topics:      make(map[string][]*Listener),
		onSubscribe: make(map[string][]pkg.PubSubSubscriptionHandler),
	}
}

func (ps *PubSub) AcceptConnection(conn pkg.PubSubConnection) {
	ps.connections = append(ps.connections, conn)
}

func (ps *PubSub) Run() {
	for msg := range ps.c {
		switch msg.operation {
		case operationPublish:
			ps.publish(msg.source, msg.topic, msg.data)
		case operationSubscribe:
			subscriptionParameters, ok := msg.data.(json.RawMessage)
			if ok {
				ps.Subscribe(msg.source, msg.topic, subscriptionParameters)
			}

		default:
			fmt.Printf("Unhandled operation: %v\n", msg.operation)
		}
	}
}

func (ps *PubSub) publish(source pkg.PubSubSource, topic string, data interface{}) {
	ps.topicsMutex.Lock()

	for _, l := range ps.topics[topic] {
		go func(listener *Listener) {
			err := listener.Publish(source, topic, data)
			if err != nil {
				ps.UnsubscribeAll(listener)
				fmt.Println(err)
			}
		}(l)
	}

	ps.topicsMutex.Unlock()
}

func (ps *PubSub) Publish(source pkg.PubSubSource, topic string, data interface{}) {
	ps.c <- Message{
		operation: operationPublish,
		topic:     topic,
		data:      data,
		source:    source,
	}
}

type pubsubMessage struct {
	Type  string
	Topic string
	Data  interface{} `json:",omitempty"`
}

type incomingMessage struct {
	Type  string
	Topic string
	Data  json.RawMessage
}

// HandleJSON handles a json blob (bytes) from the given source (source)
func (ps *PubSub) HandleJSON(source pkg.PubSubSource, bytes []byte) error {
	var msg incomingMessage
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		return err
	}

	switch msg.Type {
	case "Publish":
		ps.c <- Message{
			operation: operationPublish,
			topic:     msg.Topic,
			data:      msg.Data,
			source:    source,
		}
	case "Subscribe":
		ps.c <- Message{
			operation: operationSubscribe,
			topic:     msg.Topic,
			data:      msg.Data,
			source:    source,
		}
	}

	return nil
}

func (ps *PubSub) Subscribe(source pkg.PubSubSource, topic string, parameters json.RawMessage) {
	successfulAuthorization := ps.notifySubscriptionHandlers(source, topic, parameters)
	if !successfulAuthorization {
		fmt.Printf("[%s] Failed authorization:\n", topic)
		// fmt.Printf("[%s] Failed authorization: %+v\n", topic, auth)
		return
	}

	{
		ps.topicsMutex.Lock()
		defer ps.topicsMutex.Unlock()
		l := &Listener{source.Connection(), SubscriptionTypeContinuous}

		ps.topics[topic] = append(ps.topics[topic], l)
	}
}

func (ps *PubSub) notifySubscriptionHandlers(source pkg.PubSubSource, topic string, parameters json.RawMessage) bool {
	ps.onSubscribeMutex.Lock()
	defer ps.onSubscribeMutex.Unlock()

	for _, handler := range ps.onSubscribe[topic] {
		err, successfulAuthorization := handler.ConnectionSubscribed(source, topic, parameters)
		if err != nil {
			fmt.Println("Error in subscription handler:", err)
		}

		if !successfulAuthorization {
			return false
		}
	}

	return true
}

func (ps *PubSub) HandleSubscribe(connection pkg.PubSubSubscriptionHandler, topic string) {
	ps.onSubscribeMutex.Lock()
	defer ps.onSubscribeMutex.Unlock()

	ps.onSubscribe[topic] = append(ps.onSubscribe[topic], connection)
}

func (ps *PubSub) SubscribeOnce(source pkg.PubSubSource, topic string) {
	ps.topicsMutex.Lock()
	defer ps.topicsMutex.Unlock()

	ps.topics[topic] = append(ps.topics[topic], &Listener{source.Connection(), SubscriptionTypeOnce})
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
