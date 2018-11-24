package main

import (
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/twitch"
)

var _ pkg.StreamStore = &StreamStore{}

type StreamStore struct {
	streams map[string]*twitch.Stream
}

func NewStreamStore() *StreamStore {
	s := &StreamStore{}

	return s
}

func (s *StreamStore) GetStream(account pkg.Account) pkg.Stream {
	return nil
}
