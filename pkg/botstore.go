package pkg

type BotStore interface {
	Add(Sender)

	GetBotFromName(string) Sender
	GetBotFromID(string) Sender
	GetBotFromChannel(string) Sender

	Iterate() BotStoreIterator
}

type BotStoreIterator interface {
	Next() bool
	Value() Sender
}
