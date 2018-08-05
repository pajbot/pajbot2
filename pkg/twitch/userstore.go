package twitch

import (
	"fmt"
	"log"
	"strings"

	"github.com/dankeroni/gotwitch"
	"github.com/pajlada/pajbot2/apirequest"
	"github.com/pajlada/pajbot2/pkg"
)

var _ pkg.UserStore = &UserStore{}

type UserStore struct {
	userIDMap map[string]string
}

func NewUserStore() *UserStore {
	return &UserStore{
		userIDMap: make(map[string]string),
	}
}

func onHTTPError(statusCode int, statusMessage, errorMessage string) {
	log.Println("HTTPERROR: ", errorMessage)
}

func onInternalError(err error) {
	log.Printf("internal error: %s", err)
}

func min(a, b int) int {
	if b < a {
		return b
	}

	return a
}

func (s *UserStore) GetIDs(usernames []string) map[string]string {
	userIDs := make(map[string]string)

	remainingUsernames := []string{}
	for _, username := range usernames {
		if userID, ok := s.userIDMap[username]; ok {
			userIDs[username] = userID
		} else {
			remainingUsernames = append(remainingUsernames, username)
		}
	}

	var batch []string

	for len(remainingUsernames) > 0 {
		if len(batch) == 0 {
			batch = remainingUsernames[0:min(99, len(remainingUsernames))]
			remainingUsernames = remainingUsernames[len(batch):]
		}

		onSuccess := func(data []gotwitch.User) {
			for _, user := range data {
				userIDs[user.Login] = user.ID
				s.userIDMap[user.Login] = user.ID
			}
			batch = nil
		}

		fmt.Printf("Sending request for %v\n", batch)
		apirequest.Twitch.GetUsersByLogin(batch, onSuccess, onHTTPError, onInternalError)
	}

	return userIDs
}

func (s *UserStore) GetID(username string) string {
	username = strings.ToLower(username)

	if userID, ok := s.userIDMap[username]; ok {
		return userID
	}

	var retUserID string

	onSuccess := func(data []gotwitch.User) {
		if len(data) == 0 {
			// :(
			return
		}

		retUserID = data[0].ID
	}

	apirequest.Twitch.GetUsersByLogin([]string{username}, onSuccess, onHTTPError, onInternalError)

	return retUserID
}
