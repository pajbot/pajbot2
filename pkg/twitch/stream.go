package twitch

import (
	"sync"
	"time"

	"github.com/dankeroni/gotwitch"
	"github.com/pajlada/pajbot2/pkg"
)

var _ pkg.Stream = &Stream{}

type StreamStatus struct {
	*gotwitch.Stream

	mutex *sync.RWMutex
}

func (s *StreamStatus) Update(streamData *gotwitch.Stream) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Stream = streamData
}

func (s *StreamStatus) Live() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.Stream != nil
}

func (s *StreamStatus) StartedAt() (r time.Time) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Stream != nil {
		r = s.Stream.StartedAt
	}

	return
}

type Stream struct {
	ID string

	status StreamStatus

	needsInitialPoll bool
}

func NewTwitchStream(account pkg.Account) *Stream {
	s := &Stream{
		ID: account.ID(),

		status: StreamStatus{
			mutex: &sync.RWMutex{},
		},

		needsInitialPoll: true,
	}

	return s
}

func (s *Stream) Status() pkg.StreamStatus {
	return &s.status
}

func (s *Stream) Update(stream *gotwitch.Stream) {
	s.status.Update(stream)
}

func (s *Stream) NeedsInitialPoll() bool {
	if s.needsInitialPoll {
		s.needsInitialPoll = false
		return true
	}

	return false
}
