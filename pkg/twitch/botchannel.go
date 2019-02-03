package twitch

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/eventemitter"
	"github.com/pajlada/pajbot2/pkg/modules"
	"github.com/pajlada/pajbot2/pkg/utils"
)

var _ pkg.BotChannel = &BotChannel{}

type BotChannel struct {
	streamStore pkg.StreamStore

	ID int64

	channel User
	BotUser User

	initialized bool

	// Enabled modules
	modules      []pkg.Module
	modulesMutex sync.Mutex

	sql *sql.DB

	eventEmitter *eventemitter.EventEmitter

	bot *Bot
}

func (c *BotChannel) Channel() pkg.Channel {
	return &c.channel
}

func (c *BotChannel) Bot() pkg.Sender {
	return c.bot
}

func (c *BotChannel) Say(message string) {
	c.bot.Say(&c.channel, message)
}

func (c *BotChannel) Mention(user pkg.User, message string) {
	c.bot.Mention(&c.channel, user, message)
}

func (c *BotChannel) Timeout(user pkg.User, duration int, reason string) {
	c.bot.Timeout(&c.channel, user, duration, reason)
}

func (c *BotChannel) DatabaseID() int64 {
	return c.ID
}

func (c *BotChannel) Events() *eventemitter.EventEmitter {
	return c.eventEmitter
}

func (c *BotChannel) ChannelID() string {
	return c.channel.ID()
}

func (c *BotChannel) ChannelName() string {
	return c.channel.Name()
}

func (c *BotChannel) Stream() pkg.Stream {
	return c.streamStore.GetStream(&c.channel)
}

// We assume that modulesMutex is locked already
func (c *BotChannel) sortModules() {
	sort.Slice(c.modules, func(i, j int) bool {
		return c.modules[i].Spec().Priority() < c.modules[j].Spec().Priority()
	})
}

func (c *BotChannel) getSettingsForModule(moduleID string) ([]byte, error) {
	const queryF = `
SELECT
	settings
FROM
	BotChannelModule
WHERE
	bot_channel_id=? AND module_id=?`

	row := c.sql.QueryRow(queryF, c.DatabaseID(), moduleID)

	var s sql.NullString
	err := row.Scan(&s)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return []byte(s.String), nil
}

// We assume that modulesMutex is locked already
func (c *BotChannel) enableModule(spec pkg.ModuleSpec, settings []byte) error {
	module := spec.Maker()()
	err := module.Initialize(c, settings)
	if err != nil {
		return errors.New(fmt.Sprintf("Error loading module '%s': %s\n", spec.ID(), err.Error()))
	}

	c.modules = append(c.modules, module)

	c.sortModules()

	return nil
}

func (c *BotChannel) setModuleEnabledState(moduleID string, state *bool) error {
	const queryF = `
INSERT INTO
	BotChannelModule
	(bot_channel_id, module_id, enabled)
	VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE enabled=?`

	_, err := c.sql.Exec(queryF, c.DatabaseID(), moduleID, state, state)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// We assume that modulesMutex is locked already
func (c *BotChannel) EnableModule(moduleID string) error {
	moduleID = strings.ToLower(moduleID)

	spec, ok := modules.GetModule(moduleID)
	if !ok {
		return errors.New("invalid module id")
	}

	// Check if module is enabled already

	for _, m := range c.modules {
		if m.Spec().ID() == moduleID {
			return errors.New("module already enabled")
		}
	}

	// Save enabled state
	if err := c.setModuleEnabledState(moduleID, utils.BoolPtr(true)); err != nil {
		return err
	}

	settings, err := c.getSettingsForModule(moduleID)
	if err != nil {
		return err
	}

	return c.enableModule(spec, settings)
}

// We assume that modulesMutex is locked already
func (c *BotChannel) DisableModule(moduleID string) error {
	moduleID = strings.ToLower(moduleID)

	_, ok := modules.GetModule(moduleID)
	if !ok {
		return errors.New("invalid module id")
	}

	for i, m := range c.modules {
		if m.Spec().ID() == moduleID {
			m.Disable()
			c.modules = append(c.modules[:i], c.modules[i+1:]...)

			// Save disabled state
			if err := c.setModuleEnabledState(moduleID, utils.BoolPtr(false)); err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("module isn't enabled")
}

func (c *BotChannel) Initialize(b *Bot) error {
	if c.initialized {
		return errors.New("bot channel is already initialized")
	}

	c.bot = b
	c.sql = b.sql
	c.streamStore = b.streamStore

	c.initialized = true

	c.eventEmitter = eventemitter.New()

	c.loadModules()

	c.eventEmitter.Emit("on_join", nil)

	return nil
}

type moduleConfig struct {
	DatabaseID int64
	ModuleID   string
	Enabled    sql.NullBool

	// json string
	Settings string
}

func (c *BotChannel) loadAllModuleConfigs() ([]*moduleConfig, error) {
	const queryF = `SELECT id, module_id, enabled, settings FROM BotChannelModule WHERE bot_channel_id=?`

	rows, err := c.sql.Query(queryF, c.DatabaseID())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var moduleConfigs []*moduleConfig

	for rows.Next() {
		var mc moduleConfig
		var s sql.NullString
		if err = rows.Scan(&mc.DatabaseID, &mc.ModuleID, &mc.Enabled, &s); err != nil {
			return nil, err
		}

		if s.Valid {
			mc.Settings = s.String
		}

		moduleConfigs = append(moduleConfigs, &mc)
	}

	return moduleConfigs, nil
}

func (c *BotChannel) loadModules() {
	moduleConfigs, err := c.loadAllModuleConfigs()
	if err != nil {
		panic(err)
	}

	availableModules := modules.Modules()

	c.modulesMutex.Lock()
	defer c.modulesMutex.Unlock()

	for _, spec := range availableModules {
		enabled := spec.EnabledByDefault()
		var settings []byte

		var cfg *moduleConfig

		for _, moduleConfig := range moduleConfigs {
			if moduleConfig.ModuleID == spec.ID() {
				cfg = moduleConfig
				break
			}
		}

		if cfg != nil {
			if cfg.Enabled.Valid {
				enabled = cfg.Enabled.Bool
			}

			settings = []byte(cfg.Settings)
		}

		if enabled {
			c.enableModule(spec, settings)
		}
		// Fetch config for this module from SQL
	}
}

func (c *BotChannel) onModules(cb func(module pkg.Module) error) (err error) {
	c.modulesMutex.Lock()
	defer c.modulesMutex.Unlock()

	for _, module := range c.modules {
		if err = cb(module); err != nil {
			return
		}
	}

	return
}

func (c *BotChannel) handleMessage(user pkg.User, message *TwitchMessage, action pkg.Action) error {
	c.eventEmitter.Emit("on_msg", nil)

	return c.onModules(func(module pkg.Module) error {
		return module.OnMessage(c, user, message, action)
	})
}

func (c *BotChannel) handleWhisper(user pkg.User, message *TwitchMessage) error {
	return c.onModules(func(module pkg.Module) error {
		return module.OnWhisper(c, user, message)
	})
}
