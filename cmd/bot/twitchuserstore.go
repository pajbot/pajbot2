package main

import (
	"strings"
	"sync"

	"github.com/dankeroni/gotwitch"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/apirequest"
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

func min(a, b int) int {
	if b < a {
		return b
	}

	return a
}

func (s *UserStore) GetIDs(names []string) (ids map[string]string) {
	ids = make(map[string]string)

	remainingNames := []string{}
	s.idsMutex.Lock()

	for _, name := range names {
		if id, ok := s.ids[name]; ok {
			ids[name] = id
		} else {
			remainingNames = append(remainingNames, name)
		}
	}

	s.idsMutex.Unlock()

	var batch []string

	for len(remainingNames) > 0 {
		if len(batch) == 0 {
			batch = remainingNames[0:min(99, len(remainingNames))]
			remainingNames = remainingNames[len(batch):]
		}

		onSuccess := func(data []gotwitch.User) {
			s.idsMutex.Lock()
			defer s.idsMutex.Unlock()
			s.namesMutex.Lock()
			defer s.namesMutex.Unlock()

			for _, user := range data {
				ids[user.Login] = user.ID
				s.save(user.ID, user.Login)
			}
			batch = nil
		}

		apirequest.Twitch.GetUsersByLogin(batch, onSuccess, onHTTPError, onInternalError)
	}

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

	onSuccess := func(data []gotwitch.User) {
		if len(data) == 0 {
			// :(
			return
		}

		s.idsMutex.Lock()
		defer s.idsMutex.Unlock()
		s.namesMutex.Lock()
		defer s.namesMutex.Unlock()

		id = data[0].ID
		s.save(id, name)
	}

	apirequest.Twitch.GetUsersByLogin([]string{name}, onSuccess, onHTTPError, onInternalError)

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

	onSuccess := func(data []gotwitch.User) {
		if len(data) == 0 {
			// :(
			return
		}

		s.idsMutex.Lock()
		defer s.idsMutex.Unlock()
		s.namesMutex.Lock()
		defer s.namesMutex.Unlock()

		name = data[0].Login
		s.save(id, name)
	}

	apirequest.Twitch.GetUsers([]string{id}, onSuccess, onHTTPError, onInternalError)

	return
}

func (s *UserStore) save(id, name string) {
	s.names[id] = name
	s.ids[name] = id
}
