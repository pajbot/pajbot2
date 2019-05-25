package pkg

import (
	"sync"

	"github.com/pajbot/utils"
)

type UserStoreRequest struct {
	ids   map[string]bool
	names map[string]bool
}

func NewUserStoreRequest() *UserStoreRequest {
	return &UserStoreRequest{
		ids:   make(map[string]bool),
		names: make(map[string]bool),
	}
}

func (r *UserStoreRequest) AddID(id string) {
	r.ids[id] = true
}

func (r *UserStoreRequest) AddName(name string) {
	r.names[name] = true
}

// Returns two values:
// First value: map where key is a user ID pointing at a user name
// Second value: map where key is a user name pointing at a user ID
func (r *UserStoreRequest) Execute(userStore UserStore) (names map[string]string, ids map[string]string) {
	var wg sync.WaitGroup
	if len(r.ids) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			names = userStore.GetNames(utils.SBKey(r.ids))
		}()
	}

	if len(r.names) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ids = userStore.GetIDs(utils.SBKey(r.names))
		}()
	}

	wg.Wait()

	return
}

type UserStore interface {
	// Input: Lowercased twitch usernames
	// Returns: user IDs as strings in no specific order, and a bool indicating whether the user needs to exhaust the list first and wait
	GetIDs([]string) map[string]string

	GetID(string) string

	GetName(string) string

	// Input: list of twitch IDs
	// Returns: map of twitch IDs pointing at twitch usernames
	GetNames([]string) map[string]string
}
