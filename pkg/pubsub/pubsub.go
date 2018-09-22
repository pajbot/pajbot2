package pubsub

import (
	"fmt"
	"sync"
)

type topic struct {
}

type Connection interface {
	MessageReceived(topic string, bytes []byte) error
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
)

type Message struct {
	operation operationType
	topic     string
	data      interface{}
}

type PubSub struct {
	c           chan (Message)
	connections []Connection

	topicsMutex sync.Mutex
	topics      map[string][]*Listener
}

func New() *PubSub {
	return &PubSub{
		c:      make(chan (Message)),
		topics: make(map[string][]*Listener),
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
				ps.publish(msg.topic, msg.data)
			}
			fmt.Printf("Got message %#v\n", msg)
		}
	}
}

func (ps *PubSub) publish(topic string, data interface{}) {
	ps.topicsMutex.Lock()

	for _, l := range ps.topics[topic] {
		go func(listener *Listener) {
			err := listener.Publish(topic, data)
			if err != nil {
				fmt.Println(err)
			}
		}(l)
	}

	ps.topicsMutex.Unlock()
}

func (ps *PubSub) Publish(topic string, data interface{}) {
	ps.c <- Message{operationPublish, topic, data}
}

func (ps *PubSub) Subscribe(connection Connection, topic string) {
	ps.topicsMutex.Lock()
	defer ps.topicsMutex.Unlock()

	ps.topics[topic] = append(ps.topics[topic], &Listener{connection, SubscriptionTypeContinuous})
}

func (ps *PubSub) SubscribeOnce(connection Connection, topic string) {
	ps.topicsMutex.Lock()
	defer ps.topicsMutex.Unlock()

	ps.topics[topic] = append(ps.topics[topic], &Listener{connection, SubscriptionTypeOnce})
}
