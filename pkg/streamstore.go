package pkg

type StreamStore interface {
	// GetStream ensures that the given Account is being followed & polled for stream status updates
	// If the account is already being followed, the Stream that was already stored is returned.
	GetStream(Account) Stream
}
