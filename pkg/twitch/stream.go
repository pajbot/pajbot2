package twitch

import (
	"sync"
	"time"

	"github.com/dankeroni/gotwitch"
	"github.com/pajbot/pajbot2/pkg"
)

var _ pkg.Stream = &Stream{}

type StreamStatus struct {
	*gotwitch.HelixStream

	mutex *sync.RWMutex
}

func (s *StreamStatus) Update(streamData *gotwitch.HelixStream) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.HelixStream == nil && streamData != nil {
		// Fire event emitter STREAM ONLINE
		// Stream just went online
	} else if s.HelixStream != nil && streamData == nil {
		// Fire event emitter STREAM OFFLINE
		// Stream just went offline
	}

	s.HelixStream = streamData
}

func (s *StreamStatus) Live() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.HelixStream != nil
}

func (s *StreamStatus) StartedAt() (r time.Time) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.HelixStream != nil {
		r = s.HelixStream.StartedAt
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

func (s *Stream) Update(stream *gotwitch.HelixStream) {
	s.status.Update(stream)
}

func (s *Stream) NeedsInitialPoll() bool {
	if s.needsInitialPoll {
		s.needsInitialPoll = false
		return true
	}

	return false
}
