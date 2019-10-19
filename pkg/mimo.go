package pkg

// MIMO is a Many In Many Out interface
// Implementation for this exists in pkg/mimo/
type MIMO interface {
	Subscriber(channelNames ...string) chan interface{}
	Publisher(channelName string) chan interface{}
}
