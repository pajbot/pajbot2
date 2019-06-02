package mimo

import (
	"fmt"
	"sync"
)

type MIMO struct {
	subscribersMutex sync.RWMutex
	subscribers      map[string][]chan interface{}
}

func New() *MIMO {
	return &MIMO{
		subscribers: make(map[string][]chan interface{}),
	}
}

func (m *MIMO) Subscriber(channelNames ...string) (out chan interface{}) {
	out = make(chan interface{}, 1)
	m.subscribersMutex.Lock()
	for _, channelName := range channelNames {
		m.subscribers[channelName] = append(m.subscribers[channelName], out)
	}
	m.subscribersMutex.Unlock()
	return
}

func (m *MIMO) Publisher(channelName string) (in chan interface{}) {
	in = make(chan interface{})

	go func() {
		for msg := range in {
			m.subscribersMutex.RLock()
			subscribers, ok := m.subscribers[channelName]
			m.subscribersMutex.RUnlock()

			if ok {
				for _, subscriber := range subscribers {
					select {
					case subscriber <- msg:
					default:
						fmt.Println("ERROR SENDING DATA TO SUBSCRIBE. MARK FOR DELETION?")
					}
				}
			}
		}
	}()

	return in
}
