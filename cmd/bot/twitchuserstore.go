package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/apirequest"
	"github.com/pajbot/utils"
	"github.com/pkg/errors"
)

var _ pkg.UserStore = &UserStore{}

type UserStore struct {
	idsMutex *sync.Mutex
	ids      map[string]string

	namesMutex *sync.Mutex
	names      map[string]string
}

func NewUserStore() *UserStore {
	return &UserStore{
		idsMutex: &sync.Mutex{},
		ids:      make(map[string]string),

		namesMutex: &sync.Mutex{},
		names:      make(map[string]string),
	}
}

func (s *UserStore) GetIDs(names []string) (ids map[string]string) {
	ids = make(map[string]string)

	remaining := []string{}
	s.idsMutex.Lock()

	for _, name := range names {
		if id, ok := s.ids[name]; ok {
			ids[name] = id
		} else {
			remaining = append(remaining, name)
		}
	}

	s.idsMutex.Unlock()

	var wg sync.WaitGroup

	fmt.Println("GetIDs", remaining)
	batches, _ := utils.ChunkStringSlice(remaining, 100)
	for _, batch := range batches {
		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()
			data, err := apirequest.TwitchWrapper.GetUsersByLogin(batch)

			if err != nil {
				fmt.Println("API ERROR. maybe retry")
				return
			}

			s.idsMutex.Lock()
			defer s.idsMutex.Unlock()
			s.namesMutex.Lock()
			defer s.namesMutex.Unlock()

			for _, user := range data {
				ids[user.Login] = user.ID
				s.save(user.ID, user.Login)
			}
		}(batch)
	}

	wg.Wait()

	return
}

func (s *UserStore) GetID(name string) (id string) {
	var ok bool
	name = strings.ToLower(name)

	s.idsMutex.Lock()
	id, ok = s.ids[name]
	s.idsMutex.Unlock()
	if ok {
		return
	}

	fmt.Println("GetID", id)
	response, err := apirequest.TwitchWrapper.GetUsersByLogin([]string{name})
	if err != nil {
		return
	}

	if len(response) == 0 {
		return
	}

	s.idsMutex.Lock()
	defer s.idsMutex.Unlock()
	s.namesMutex.Lock()
	defer s.namesMutex.Unlock()

	id = response[0].ID
	s.save(id, name)

	return
}

func (s *UserStore) GetName(id string) (name string) {
	var ok bool

	s.namesMutex.Lock()
	name, ok = s.names[id]
	s.namesMutex.Unlock()
	if ok {
		return
	}

	fmt.Println("GetName", id)
	response, err := apirequest.TwitchWrapper.GetUsersByID([]string{id})
	if err != nil {
		return
	}

	if len(response) == 0 {
		// :(
		return
	}

	s.idsMutex.Lock()
	defer s.idsMutex.Unlock()
	s.namesMutex.Lock()
	defer s.namesMutex.Unlock()

	name = response[0].Login
	s.save(id, name)

	return
}

func (s *UserStore) GetNames(ids []string) (names map[string]string) {
	names = make(map[string]string)

	remaining := []string{}
	s.namesMutex.Lock()

	for _, id := range ids {
		if name, ok := s.names[id]; ok {
			names[id] = name
		} else {
			remaining = append(remaining, id)
		}
	}

	s.namesMutex.Unlock()

	var wg sync.WaitGroup

	fmt.Println("GetNames", remaining)
	batches, _ := utils.ChunkStringSlice(remaining, 100)
	for _, batch := range batches {
		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()
			data, err := apirequest.TwitchWrapper.GetUsersByID(batch)
			if err != nil {
				fmt.Println("API ERROR. maybe retry")
				return
			}

			s.idsMutex.Lock()
			defer s.idsMutex.Unlock()
			s.namesMutex.Lock()
			defer s.namesMutex.Unlock()

			for _, user := range data {
				names[user.ID] = user.Login
				s.save(user.ID, user.Login)
			}
		}(batch)
	}

	wg.Wait()

	return
}

func (s *UserStore) Hydrate(users []pkg.User) error {
	missingIDs := map[string]pkg.User{}
	missingLogins := map[string]pkg.User{}

	for _, user := range users {
		if user.GetName() == "" {
			if user.GetID() == "" {
				return errors.New("Cannot hydrate an empty user (no name & no id)")
			}

			missingLogins[user.GetID()] = user
		} else if user.GetID() == "" {
			if user.GetName() == "" {
				return errors.New("Cannot hydrate an empty user (no name & no id)")
			}

			missingIDs[user.GetName()] = user
		}
	}

	if len(missingIDs) > 0 {
		logins := make([]string, len(missingIDs))
		i := 0
		for login := range missingIDs {
			logins[i] = login
			i++
		}
		ids := s.GetIDs(logins)
		for name, id := range ids {
			fmt.Println("Updating ID of", name, "to", id)
			missingIDs[name].SetID(id)
		}
	}

	if len(missingLogins) > 0 {
		ids := make([]string, len(missingLogins))
		i := 0
		for login := range missingLogins {
			ids[i] = login
			i++
		}
		names := s.GetNames(ids)
		for id, name := range names {
			fmt.Println("Updating name of", id, "to", name)
			missingLogins[id].SetName(name)
		}
	}

	return nil
}

func (s *UserStore) save(id, name string) {
	s.names[id] = name
	s.ids[name] = id
}
