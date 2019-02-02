package twitch

import (
	"errors"

	"github.com/pajlada/pajbot2/pkg"
)

type User struct {
	id   string
	name string
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) GetID() string {
	return u.id
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) SetName(v string) {
	u.name = v
}

func (u *User) fillIn(userStore pkg.UserStore) error {
	if u.id == "" && u.name != "" {
		// Has name but not ID
		u.id = userStore.GetID(u.name)
		if u.id == "" {
			return errors.New("Unable to get ID")
		}

		return nil
	}

	if u.name == "" && u.id != "" {
		// Has ID but not name
		u.name = userStore.GetName(u.id)
		if u.name == "" {
			return errors.New("Unable to get Name")
		}

		return nil
	}

	if u.name != "" && u.id != "" {
		// User already has name and ID, no need to fetch anything
		return nil
	}

	return errors.New("User missing both ID and Name")
}

func (u *User) Valid() bool {
	return u.name != "" && u.id != ""
}
