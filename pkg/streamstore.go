package pkg

type StreamStore interface {
	GetStream(Account) Stream

	JoinStream(Account)
}
