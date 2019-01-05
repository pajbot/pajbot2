package pkg

type BotStore interface {
	Add(Sender)
	Get(string) Sender
	Iterate() BotStoreIterator
}

type BotStoreIterator interface {
	Next() bool
	Value() Sender
}
