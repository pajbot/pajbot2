package botstore

import "github.com/pajbot/pajbot2/pkg"

var _ pkg.BotStoreIterator = &BotStoreIterator{}

type BotStoreIterator struct {
	data []pkg.Sender

	index int
}

func (i *BotStoreIterator) Next() bool {
	i.index++

	return i.index < len(i.data)
}

func (i *BotStoreIterator) Value() pkg.Sender {
	return i.data[i.index]
}
