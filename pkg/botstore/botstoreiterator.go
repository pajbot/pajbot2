package botstore

import "github.com/pajlada/pajbot2/pkg"

var _ pkg.BotStoreIterator = &BotStoreIterator{}

type BotStoreIterator struct {
	data []pkg.Sender

	index int
}

func (i *BotStoreIterator) Next() bool {
	i.index++
	if i.index >= len(i.data) {
		return false
	}

	return true
}

func (i *BotStoreIterator) Value() pkg.Sender {
	return i.data[i.index]
}
