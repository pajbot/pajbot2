package channels

import (
	"sync"

	"github.com/pajbot/pajbot2/pkg"
)

var _ pkg.ChannelStore = &Store{}

type Store struct {
	data      map[string]pkg.Channel
	dataMutex sync.Mutex
}

func (s *Store) TwitchChannel(channelID string) (channel pkg.Channel) {
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()

	channel, _ = s.data[channelID]

	return
}

func (s *Store) RegisterTwitchChannel(channel pkg.Channel) {
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()

	s.data[channel.GetID()] = channel

}

func NewStore() *Store {
	return &Store{
		data: make(map[string]pkg.Channel),
	}
}
