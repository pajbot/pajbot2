// TODO: fix a proper api poller or something
// I need an easier way to say "hey I need 500 stream statuses", and it should then call it 5 times with 100 users in each request
// Check out API rate limiting
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/apirequest"
	"github.com/pajbot/pajbot2/pkg/twitch"
	"github.com/pajbot/utils"
)

const PollInterval = 5 * time.Second

var _ pkg.StreamStore = &StreamStore{}

type StreamStore struct {
	mutex   *sync.Mutex
	streams map[string]*twitch.Stream
}

func NewStreamStore() *StreamStore {
	s := &StreamStore{
		mutex:   &sync.Mutex{},
		streams: make(map[string]*twitch.Stream),
	}

	return s
}

func (s *StreamStore) PollStreams() {
	s.mutex.Lock()

	var remaining []string

	for streamID, stream := range s.streams {
		if stream.NeedsInitialPoll() {
			remaining = append(remaining, streamID)
		}
	}

	s.mutex.Unlock()

	if len(remaining) == 0 {
		return
	}

	fmt.Println("Poll streams")

	var wg sync.WaitGroup

	batches, _ := utils.ChunkStringSlice(remaining, 100)
	for _, batch := range batches {
		wg.Add(1)

		// Subscribe to the streams webhook topic for the channels as well
		go func(batch []string) {
			for _, userID := range batch {
				go func(userID string) {
					apirequest.TwitchWrapper.EventSubSubscribe(helix.EventSubTypeStreamOnline, userID)
					apirequest.TwitchWrapper.EventSubSubscribe(helix.EventSubTypeStreamOffline, userID)
				}(userID)
			}
		}(batch)

		go func(batch []string) {
			defer wg.Done()
			data, err := apirequest.TwitchWrapper.GetStreams(batch, nil)
			if err != nil {
				fmt.Println("api error:", err)
			}

			s.mutex.Lock()
			defer s.mutex.Unlock()
			for _, activeStream := range data {
				if stream, ok := s.streams[activeStream.UserID]; ok {
					stream.Update(&activeStream)
				}
			}
		}(batch)
	}

	wg.Wait()
}

func (s *StreamStore) Run() {
	for range time.After(PollInterval) {
		s.PollStreams()
	}
}

func (s *StreamStore) GetStream(account pkg.Account) pkg.Stream {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if stream, ok := s.streams[account.ID()]; ok {
		return stream
	}

	stream := twitch.NewTwitchStream(account)

	s.streams[account.ID()] = stream

	return stream
}
