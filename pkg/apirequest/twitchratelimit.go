package apirequest

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type TwitchRateLimit struct {
	mutex *sync.RWMutex

	// The rate at which points are added to your bucket. This is the average number of requests per minute you can make over an extended period of time.
	Limit int

	// The number of points you have left to use.
	Remaining int

	// A timestamp of when your bucket is reset to full.
	Reset time.Time
}

func NewTwitchRateLimit() TwitchRateLimit {
	return TwitchRateLimit{
		mutex: &sync.RWMutex{},
	}
}

func (l *TwitchRateLimit) Update(r *http.Response) {
	limit := r.Header.Get("Ratelimit-Limit")
	remaining := r.Header.Get("Ratelimit-Remaining")
	reset := r.Header.Get("Ratelimit-Reset")

	if limit == "" || remaining == "" || reset == "" {
		return
	}

	nLimit, err := strconv.Atoi(limit)
	if err != nil {
		fmt.Println("Error parsing limit from", limit)
	}
	nRemaining, err := strconv.Atoi(remaining)
	if err != nil {
		fmt.Println("Error parsing remaining from", remaining)
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.Limit = nLimit
	l.Remaining = nRemaining
}

func (l *TwitchRateLimit) String() string {
	return fmt.Sprintf("[RateLimit Limit=%d Remaining=%d Reset=%s]", l.Limit, l.Remaining, l.Reset)
}
