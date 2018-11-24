package pkg

type PubSub interface {
	Subscribe(source PubSubSource, topic string)
	Publish(source PubSubSource, topic string, data interface{})

	HandleSubscribe(connection PubSubSubscriptionHandler, topic string)
	HandleJSON(source PubSubSource, bytes []byte) error
}

// PubSubConnection is an interface where a JSON message can be written to
type PubSubConnection interface {
	MessageReceived(source PubSubSource, topic string, bytes []byte) error
}

type PubSubSubscriptionHandler interface {
	ConnectionSubscribed(source PubSubSource, topic string) (error, bool)
}

// PubSubSource is an interface that is responsible for a message being written into pubsub
// This will be responsible for checking authorization
type PubSubSource interface {
	IsApplication() bool
	Connection() PubSubConnection
	AuthenticatedUser() User
}

type PubSubBan struct {
	Channel string
	Target  string
	Reason  string
}

type PubSubTimeout struct {
	Channel  string
	Target   string
	Reason   string
	Duration uint32
}

type PubSubUntimeout struct {
	Channel string
	Target  string
}

type PubSubUser struct {
	ID   string
	Name string
}

type PubSubBanEvent struct {
	Channel PubSubUser
	Target  PubSubUser
	Source  PubSubUser
	Reason  string
}

type PubSubTimeoutEvent struct {
	Channel  PubSubUser
	Target   PubSubUser
	Source   PubSubUser
	Duration int
	Reason   string
}
