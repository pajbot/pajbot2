package pubsub

import (
	"encoding/json"
	"fmt"

	"github.com/pajlada/pajbot2/pkg"
)

type Listener struct {
	connection       pkg.PubSubConnection
	subscriptionType SubscriptionType
}

func (l *Listener) Publish(source pkg.PubSubSource, topic string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Unable to unmarshal %#v\n", data)
		// we don't treat this as an actual error
		return nil
	}

	err = l.connection.MessageReceived(source, topic, bytes)
	if err != nil {
		return err
	}

	return nil
}
