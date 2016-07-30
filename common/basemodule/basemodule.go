package basemodule

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
	"github.com/pajlada/pajbot2/helper"
	"github.com/pajlada/pajbot2/plog"
	"github.com/pajlada/pajbot2/redismanager"
)

var log = plog.GetLogger()

// BaseModule includes information on whether it's enabled or not
type BaseModule struct {
	// Unique identifier of the module
	ID string

	// Whether the module is enabled or not
	// Valid values:
	// null = default module value
	// true = enabled
	// false = disabled
	Enabled *bool `json:"enabled"`

	// Whether the module is enabled by default or not
	EnabledDefault bool

	// Level required to call the module
	// Valid values:
	// null = default level required as set in the module
	// 0-2000 = int value, level required to call module
	LevelRequired *int64 `json:"level_required"`

	LevelRequiredDefault int

	// Level required to bypass the module
	// Valid values:
	// null = default level required as set in the module
	// -1 = Unbypassable
	// 0-2000 = int value, level required to bypass
	LevelBypass *int64 `json:"level_bypass"`

	LevelBypassDefault int
}

// SetDefaults sets the defaults values on the given module.
func (m *BaseModule) SetDefaults(id string) {
	m.EnabledDefault = false
	m.LevelRequiredDefault = 100
	m.LevelBypassDefault = -1
	m.ID = id
}

// GetState xD
func (m *BaseModule) GetState() *BaseModule {
	return m
}

// SetEnabled xD
func (m *BaseModule) SetEnabled() {
	m.Enabled = helper.GetTrueP()
}

// SetDisabled xD
func (m *BaseModule) SetDisabled() {
	m.Enabled = helper.GetFalseP()
}

// IsEnabled xD
func (m *BaseModule) IsEnabled() bool {
	if m.Enabled == nil {
		return m.EnabledDefault
	}
	return *m.Enabled == true
}

// GetLevelRequired xD
func (m *BaseModule) GetLevelRequired() int {
	if m.LevelRequired == nil {
		return m.LevelRequiredDefault
	}
	return int(*m.LevelRequired)
}

// GetLevelBypass xD
func (m *BaseModule) GetLevelBypass() int {
	if m.LevelBypass == nil {
		return m.LevelBypassDefault
	}
	return int(*m.LevelBypass)
}

// FetchSettings xD
func (m *BaseModule) FetchSettings(r *redismanager.RedisManager, channelName string) []byte {
	conn := r.Pool.Get()
	defer conn.Close()
	data, err := redis.Bytes(conn.Do("HGET", channelName+":modules:settings", m.ID))
	if err != nil {
		// No settings, or the connection is fucked
		// Either way, no data to return
		return nil
	}

	return data
}

// ParseState xD
func (m *BaseModule) ParseState(r *redismanager.RedisManager, channelName string) {
	conn := r.Pool.Get()
	defer conn.Close()
	data, err := redis.Bytes(conn.Do("HGET", channelName+":modules:state", m.ID))
	if err != nil {
		return
	}

	err = json.Unmarshal(data, m)
	if err != nil {
		log.Error(err)
		return
	}
}

// ParseSettings xD
func (m *BaseModule) ParseSettings(r *redismanager.RedisManager, channelName string) {
	data := m.FetchSettings(r, channelName)

	if data == nil {
		return
	}

	err := json.Unmarshal(data, m)
	if err != nil {
		log.Error(err)
		return
	}
}
