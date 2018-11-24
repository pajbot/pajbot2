package pkg

type StreamStore interface {
	GetStream(Account) Stream
}
