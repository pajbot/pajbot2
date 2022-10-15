package pkg

import (
	"time"

	"github.com/nicklaw5/helix/v2"
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
