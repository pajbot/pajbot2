package pkg

import (
	"time"

	"github.com/nicklaw5/helix"
)

type StreamStatus interface {
	Live() bool
	StartedAt() time.Time
}

type Stream interface {
	Status() StreamStatus

	// Update forwards the given helix data to its internal status
	Update(*helix.Stream)
}
