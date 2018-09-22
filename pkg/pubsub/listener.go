package pubsub

import (
	"encoding/json"
	"fmt"
)

type Listener struct {
	connection       Connection
	subscriptionType SubscriptionType
}

func (l *Listener) Publish(topic string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Unable to unmarshal %#v\n", data)
		// we don't treat this as an actual error
		return nil
	}

	err = l.connection.MessageReceived(topic, bytes)
	if err != nil {
		return err
	}

	return nil
}
