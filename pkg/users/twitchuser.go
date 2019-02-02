package users

import (
	"database/sql"
	"errors"
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

	channelPermissions map[string]*permissionSet
}

var _ pkg.User = &TwitchUser{}

func NewTwitchUser(user twitch.User, id string) *TwitchUser {
	return &TwitchUser{
		User: user,

		ID: id,
	}
}

func NewSimpleTwitchUser(userID, userName string) *TwitchUser {
	u := &TwitchUser{
		User: twitch.User{
			UserID:      userID,
			Username:    userName,
			DisplayName: userName,
		},

		ID: userID,
	}

	return u
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

func (u *TwitchUser) loadPermissions() error {
	p, err := GetUserGlobalPermissions(u.GetID())
	if err != nil {
		return err
	}

	u.permissionsLoaded = true
	u.permissions = p

	return nil
}

func (u *TwitchUser) HasPermission(channel pkg.Channel, permission pkg.Permission) bool {
	if u.HasChannelPermission(channel, permission) {
		return true
	}

	return u.HasGlobalPermission(permission)
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
	return u.GetName() == channel.GetName()
}

func (u TwitchUser) GetBadges() map[string]int {
	return u.Badges
}

func GetUserPermissions(userID, channelID string) (pkg.Permission, error) {
	switch channelID {
	case "global":
		return GetUserGlobalPermissions(userID)
	default:
		return GetUserChannelPermissions(userID, channelID)
	}
}

func SetUserPermissions(userID, channelID string, newPermissions pkg.Permission) error {
	switch channelID {
	case "global":
		return SetUserGlobalPermissions(userID, newPermissions)
	default:
		return SetUserChannelPermissions(userID, channelID, newPermissions)
	}
}

func GetUserGlobalPermissions(userID string) (pkg.Permission, error) {
	var permissions pkg.Permission

	if userID == "" {
		return permissions, errors.New("missing user id or channel id")
	}

	const queryF = "SELECT permissions FROM `TwitchUserGlobalPermission` WHERE `twitch_user_id`=?;"

	var permissionsBytes []uint8
	err := _server.sql.QueryRow(queryF, userID).Scan(&permissionsBytes)
	if err != nil {
		if err == sql.ErrNoRows {
			return permissions, nil
		}

		return permissions, err
	}

	permissions = pkg.Permission(utils.BytesToUint64(permissionsBytes))

	return permissions, nil
}

func GetUserChannelPermissions(userID, channelID string) (pkg.Permission, error) {
	var permissions pkg.Permission

	if userID == "" || channelID == "" {
		return permissions, errors.New("missing user id or channel id")
	}

	const queryF = "SELECT permissions FROM `TwitchUserChannelPermission` WHERE `twitch_user_id`=? AND `channel_id`=?;"

	rows, err := _server.sql.Query(queryF, userID, channelID)
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

	const queryF = "INSERT INTO TwitchUserChannelPermission (twitch_user_id, channel_id, permissions) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE permissions=?;"

	permissionBytes := utils.Uint64ToBytes(uint64(permission))

	rows, err := _server.sql.Query(queryF, userID, channelID, permissionBytes, permissionBytes)
	if err != nil {
		return err
	}

	err = rows.Close()
	if err != nil {
		return err
	}

	return nil
}

func SetUserGlobalPermissions(userID string, permission pkg.Permission) error {
	const queryF = `
INSERT INTO TwitchUserGlobalPermission
	(twitch_user_id, permissions)
	VALUES (?, ?)
	ON DUPLICATE KEY UPDATE permissions=?;
	`

	if userID == "" {
		return errors.New("missing user id or channel id")
	}

	permissionBytes := utils.Uint64ToBytes(uint64(permission))

	_, err := _server.sql.Exec(queryF, userID, permissionBytes, permissionBytes)
	if err != nil {
		return err
	}

	return nil
}

func HasGlobalPermission(userID string, permission pkg.Permission) (hasPermission bool, err error) {
	globalPermissions, err := GetUserGlobalPermissions(userID)
	if err != nil {
		return
	}

	hasPermission = (globalPermissions & permission) != 0
	return
}

func HasChannelPermission(userID, channelID string, permission pkg.Permission) (hasPermission bool, err error) {
	globalPermissions, err := GetUserGlobalPermissions(userID)
	if err != nil {
		return
	}
	channelPermissions, err := GetUserChannelPermissions(userID, channelID)
	if err != nil {
		return
	}

	hasPermission = ((globalPermissions | channelPermissions) & permission) != 0
	return
}
