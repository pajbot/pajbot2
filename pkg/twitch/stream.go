package twitch

import (
	"fmt"
	"sync"
	"time"

	"github.com/nicklaw5/helix"
	"github.com/pajbot/pajbot2/pkg"
)

var _ pkg.Stream = &Stream{}

type StreamStatus struct {
	*helix.Stream

	mutex *sync.RWMutex
}

func (s *StreamStatus) Update(streamData *helix.Stream) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Stream == nil && streamData != nil {
		fmt.Println("Stream went online")
		s.Stream = streamData
		// Fire event emitter STREAM ONLINE
		// Stream just went online
	} else if s.Stream != nil && streamData == nil {
		fmt.Println("Stream went offline")
		s.Stream = streamData
		// Fire event emitter STREAM OFFLINE
		// Stream just went offline
	} else if s.Stream != nil && streamData != nil {
		// We have stream data already, but the incoming values are also set. Figure out which to use
		// TODO: MERGE STREAM
		if streamData.ID != "" {
			s.Stream.ID = streamData.ID
		}
		if streamData.UserID != "" {
			s.Stream.UserID = streamData.UserID
		}
		if streamData.UserLogin != "" {
			s.Stream.UserLogin = streamData.UserLogin
		}
		if streamData.UserName != "" {
			s.Stream.UserName = streamData.UserName
		}
		if streamData.GameID != "" {
			s.Stream.GameID = streamData.GameID
		}
		if streamData.GameName != "" {
			s.Stream.GameName = streamData.GameName
		}
		if len(streamData.TagIDs) != 0 {
			s.Stream.TagIDs = streamData.TagIDs
		}
		if streamData.IsMature {
			s.Stream.IsMature = streamData.IsMature
		}
		if streamData.Type != "" {
			s.Stream.Type = streamData.Type
		}
		if streamData.Title != "" {
			s.Stream.Title = streamData.Title
		}
		if streamData.ViewerCount != 0 {
			s.Stream.ViewerCount = streamData.ViewerCount
		}
		s.Stream.StartedAt = streamData.StartedAt
		if streamData.Language != "" {
			s.Stream.Language = streamData.Language
		}
		if streamData.ThumbnailURL != "" {
			s.Stream.ThumbnailURL = streamData.ThumbnailURL
		}
	} else {
		// Our own stream was nil while streamData is not nil. Just replace our data
		s.Stream = streamData
	}
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

func (s *Stream) Update(stream *helix.Stream) {
	s.status.Update(stream)
}

func (s *Stream) NeedsInitialPoll() bool {
	if s.needsInitialPoll {
		s.needsInitialPoll = false
		return true
	}

	return false
}
