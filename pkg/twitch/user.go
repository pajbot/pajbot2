package twitch

import (
	"errors"

	"github.com/pajlada/pajbot2/pkg"
)

type User struct {
	ID   string
	Name string
}

func (u *User) fillIn(userStore pkg.UserStore) error {
	if u.ID == "" && u.Name != "" {
		// Has name but not ID
		u.ID = userStore.GetID(u.Name)
		if u.ID == "" {
			return errors.New("Unable to get ID")
		}

		return nil
	}

	if u.Name == "" && u.ID != "" {
		// Has ID but not name
		u.Name = userStore.GetName(u.ID)
		if u.Name == "" {
			return errors.New("Unable to get Name")
		}

		return nil
	}

	if u.Name != "" && u.ID != "" {
		// User already has name and ID, no need to fetch anything
		return nil
	}

	return errors.New("User missing both ID and Name")
}

func (u *User) Valid() bool {
	return u.Name != "" && u.ID != ""
}
