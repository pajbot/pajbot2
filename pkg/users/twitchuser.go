package users

import (
	"fmt"
	"log"

	twitch "github.com/gempir/go-twitch-irc"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/utils"
)

type TwitchUser struct {
	twitch.User

	ID string

	permissionsLoaded bool
	permissions       pkg.Permission
}

func (u *TwitchUser) loadPermissions() error {
	const queryF = "SELECT permissions FROM `twitch_user_permissions` WHERE `twitch_user_id`=?;"

	u.permissionsLoaded = true

	rows, err := _server.sql.Session.Query(queryF, u.GetID())
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		var permissionsBytes []uint8
		err := rows.Scan(&permissionsBytes)
		if err != nil {
			return err
		}

		fmt.Printf("xd %#v\n", permissionsBytes)

		u.permissions = pkg.Permission(utils.BytesToUint64(permissionsBytes))

		// u.permissions = uint64(permissionsBytes)

		fmt.Printf("user permissions: %#v\n", u.permissions)
	}

	return nil
}

func (u *TwitchUser) HasPermission(permission pkg.Permission) bool {
	if !u.permissionsLoaded {
		err := u.loadPermissions()
		if err != nil {
			log.Println("Error loading permissions:", err)
		}
	}

	return (u.permissions & permission) != 0
}

func (u TwitchUser) GetName() string {
	return u.Username

}
func (u TwitchUser) GetDisplayName() string {
	return ""
}
func (u TwitchUser) GetID() string {
	return u.ID
}

func (u TwitchUser) IsModerator() bool {
	return u.UserType == "mod"
}

func (u TwitchUser) IsBroadcaster(channel pkg.Channel) bool {
	if channel == nil {
		return false
	}

	// TODO: tolower?
	return u.GetName() == channel.GetChannel()
}
