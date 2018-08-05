package users

import (
	"errors"
	"fmt"
	"log"

	twitch "github.com/gempir/go-twitch-irc"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/utils"
)

var _ pkg.User = &TwitchUser{}

func NewTwitchUser(user twitch.User, id string) *TwitchUser {
	return &TwitchUser{
		User: user,

		ID: id,
	}
}

type permissionSet struct {
	loaded      bool
	permissions pkg.Permission
}

func (p *permissionSet) load(channelID, userID string) error {
	if !p.loaded {
		var err error
		p.permissions, err = GetUserChannelPermissions(userID, channelID)
		if err != nil {
			return err
		}

		p.loaded = true

		return nil
	}

	return nil
}

type TwitchUser struct {
	twitch.User

	ID string

	permissionsLoaded bool
	permissions       pkg.Permission

	channelPermissions map[string]*permissionSet
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

func (u *TwitchUser) HasGlobalPermission(permission pkg.Permission) bool {
	if !u.permissionsLoaded {
		err := u.loadPermissions()
		if err != nil {
			log.Println("Error loading permissions:", err)
		}
	}

	return (u.permissions & permission) != 0
}

func (u *TwitchUser) HasChannelPermission(channel pkg.Channel, permission pkg.Permission) bool {
	if u.channelPermissions == nil {
		u.channelPermissions = make(map[string]*permissionSet)
	}

	channelPermissionSet := u.channelPermissions[channel.GetID()]
	if channelPermissionSet == nil {
		channelPermissionSet = &permissionSet{}
		u.channelPermissions[channel.GetID()] = channelPermissionSet
	}

	err := channelPermissionSet.load(channel.GetID(), u.GetID())
	if err != nil {
		log.Println("Error loading permissions:", err)
	}

	return (channelPermissionSet.permissions & permission) != 0
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

func GetUserChannelPermissions(userID, channelID string) (pkg.Permission, error) {
	var permissions pkg.Permission

	if userID == "" || channelID == "" {
		return permissions, errors.New("missing user id or channel id")
	}

	const queryF = "SELECT permissions FROM `twitch_user_channel_permissions` WHERE `twitch_user_id`=? AND `channel_id`=?;"

	rows, err := _server.sql.Session.Query(queryF, userID, channelID)
	if err != nil {
		return permissions, err
	}
	defer rows.Close()

	if rows.Next() {
		var permissionsBytes []uint8
		err := rows.Scan(&permissionsBytes)
		if err != nil {
			return permissions, err
		}

		permissions = pkg.Permission(utils.BytesToUint64(permissionsBytes))
	}

	return permissions, nil

}

func SetUserChannelPermissions(userID, channelID string, permission pkg.Permission) error {
	if userID == "" || channelID == "" {
		return errors.New("missing user id or channel id")
	}

	const queryF = "INSERT INTO twitch_user_channel_permissions (twitch_user_id, channel_id, permissions) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE permissions=?;"

	permissionBytes := utils.Uint64ToBytes(uint64(permission))

	rows, err := _server.sql.Session.Query(queryF, userID, channelID, permissionBytes, permissionBytes)
	if err != nil {
		return err
	}

	err = rows.Close()
	if err != nil {
		return err
	}

	return nil
}
