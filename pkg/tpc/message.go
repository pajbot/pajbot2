package tpc

import "encoding/json"

type incomingMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// only really relevant for our sending of messages
type outgoingMessage struct {
	Type string   `json:"type"`
	Data []uint64 `json:"data"`
}
