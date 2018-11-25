package pkg

import "time"

type StreamStatus interface {
	Live() bool
	StartedAt() time.Time
}

type Stream interface {
	Status() StreamStatus
}
