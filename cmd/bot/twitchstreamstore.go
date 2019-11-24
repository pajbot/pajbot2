// TODO: fix a proper api poller or something
// I need an easier way to say "hey I need 500 stream statuses", and it should then call it 5 times with 100 users in each request
// Check out API rate limiting
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/dankeroni/gotwitch"
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
					err := apirequest.TwitchWrapper.WebhookSubscribe(gotwitch.WebhookTopicStreams, userID)
					if err != nil {
						fmt.Println("Error subscribing to webhook for user", userID, err)
					}
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

func (s *StreamStore) JoinStream(account pkg.Account) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.streams[account.ID()]; !ok {
		// Insert stream
		s.streams[account.ID()] = twitch.NewTwitchStream(account)
	}
}
