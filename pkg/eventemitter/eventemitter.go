package eventemitter

import (
	"errors"
	"sync"
)

// type EventEmitter interface {
// 	Listen(event string, cb interface{}, priority int) (*Listener, error)
// 	Emit(event string, arguments map[string]interface{})
// }

var (
	ErrBadCallback = errors.New("Bad callback passed through to Listen")
)

type Listener struct {
	Disconnected bool

	cb       interface{}
	priority int
}

type EventEmitter struct {
	listenersMutex sync.Mutex
	listeners      map[string][]*Listener
}

func New() *EventEmitter {
	return &EventEmitter{
		listenersMutex: sync.Mutex{},
		listeners:      make(map[string][]*Listener),
	}
}

func (e *EventEmitter) Listen(event string, cb interface{}, priority int) (*Listener, error) {
	switch cb.(type) {
	case func(map[string]interface{}) error:
	case func() error:
	default:
		return nil, ErrBadCallback
	}
	e.listenersMutex.Lock()
	defer e.listenersMutex.Unlock()

	l := &Listener{
		cb:       cb,
		priority: priority,
	}

	// TODO: sort by priority
	e.listeners[event] = append(e.listeners[event], l)

	return l, nil
}

func (e *EventEmitter) Emit(event string, arguments map[string]interface{}) (n int, err error) {
	e.listenersMutex.Lock()
	defer e.listenersMutex.Unlock()

	listeners, ok := e.listeners[event]
	if !ok {
		return
	}

	for _, listener := range listeners {
		if listener.Disconnected {
			// TODO: Remove from listeners
			continue
		}

		switch cb := listener.cb.(type) {
		case func(map[string]interface{}) error:
			err = cb(arguments)
			if err != nil {
				return
			}
			n++
		case func() error:
			err = cb()
			if err != nil {
				return
			}
			n++
		}
	}

	return
}
