package botstore

import (
	"github.com/pajlada/pajbot2/pkg"
	"strings"
)

var _ pkg.BotStore = &BotStore{}

type BotStore struct {
	store []pkg.Sender
}

func New() *BotStore {
	return &BotStore{}
}

func (s *BotStore) Add(bot pkg.Sender) {
	s.store = append(s.store, bot)
}

func (s *BotStore) Get(name string) pkg.Sender {
	for _, b := range s.store {
		if b.TwitchAccount().Name() == strings.ToLower(name) {
			return b
		}
	}

	return nil
}

func (s *BotStore) Iterate() pkg.BotStoreIterator {
	return &BotStoreIterator{
		data:  s.store,
		index: -1,
	}
}
