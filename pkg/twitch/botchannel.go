package twitch

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/eventemitter"
	"github.com/pajbot/pajbot2/pkg/modules"
	"github.com/pajbot/utils"
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

func (c *BotChannel) Ban(user pkg.User, reason string) {
	c.bot.Ban(&c.channel, user, reason)
}

func (c *BotChannel) SingleTimeout(user pkg.User, duration int, reason string) {
	c.bot.SingleTimeout(&c.channel, user, duration, reason)
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
		return c.modules[i].Priority() < c.modules[j].Priority()
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
	// This will call the modules maker
	module := spec.Create(c)

	log.Println("Load module", spec.Name())

	// Load the modules setting (generally handled by the modules base)
	err := module.LoadSettings(settings)
	if err != nil {
		return fmt.Errorf("error loading module '%s': %s", spec.ID(), err.Error())
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
	log.Println("Enable module!!!!!!!", moduleID)
	moduleID = strings.ToLower(moduleID)

	spec, ok := modules.GetModuleSpec(moduleID)
	if !ok {
		return errors.New("invalid module id")
	}

	// Check if module is enabled already

	for _, m := range c.modules {
		if m.ID() == moduleID {
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

	for i, m := range c.modules {
		if m.ID() == moduleID {
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
	log.Println("LOAD MODULES")
	log.Println("LOAD MODULES")
	log.Println("LOAD MODULES")
	log.Println("LOAD MODULES")
	log.Println("LOAD MODULES")
	log.Println("LOAD MODULES")
	log.Println("LOAD MODULES")
	log.Println("LOAD MODULES")
	moduleConfigs, err := c.loadAllModuleConfigs()
	if err != nil {
		panic(err)
	}

	availableModules := modules.Modules()
	log.Println("Available modules:", availableModules)

	c.modulesMutex.Lock()
	defer c.modulesMutex.Unlock()

	for _, spec := range availableModules {
		log.Println("Available module:", spec)
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

func (c *BotChannel) OnModules(cb func(module pkg.Module) pkg.Actions) (actions []pkg.Actions) {
	c.modulesMutex.Lock()
	defer c.modulesMutex.Unlock()

	for _, module := range c.modules {
		// TODO: This could potentially be run in steps now. maybe all modules with same priority are run together?
		if moduleActions := cb(module); moduleActions != nil {
			actions = append(actions, moduleActions)
		}
	}

	return
}

func (c *BotChannel) resolveActions(actions []pkg.Actions) error {
	// TODO: Resolve actions smarter
	for _, action := range actions {
		for _, mute := range action.Mutes() {
			switch mute.Type() {
			case pkg.MuteTypeTemporary:
				c.Timeout(mute.User(), int(mute.Duration().Seconds()), mute.Reason())
			case pkg.MuteTypePermanent:
				c.Ban(mute.User(), mute.Reason())
			}
		}

		for _, message := range action.Messages() {
			c.Say(message.Evaluate())
		}

		for _, whisper := range action.Whispers() {
			c.bot.Whisper(whisper.User(), whisper.Content())
		}

	}

	fmt.Println("Got actions:", actions)

	return nil
}

func (c *BotChannel) HandleMessage(user pkg.User, message pkg.Message) error {
	c.eventEmitter.Emit("on_msg", nil)

	log.Println("Got message:", message.GetText())

	event := pkg.MessageEvent{
		BaseEvent: pkg.BaseEvent{
			UserStore: c.bot.GetUserStore(),
		},
		User:    user,
		Message: message,
		Channel: c.Channel(),
	}

	actions := c.OnModules(func(module pkg.Module) pkg.Actions {
		return module.OnMessage(event)
	})

	// TODO: Resolve actions smarter
	return c.resolveActions(actions)
}

func (c *BotChannel) handleWhisper(user pkg.User, message *TwitchMessage) error {
	event := pkg.MessageEvent{
		BaseEvent: pkg.BaseEvent{
			UserStore: c.bot.GetUserStore(),
		},
		User:    user,
		Message: message,
		Channel: c.Channel(),
	}

	actions := c.OnModules(func(module pkg.Module) pkg.Actions {
		return module.OnWhisper(event)
	})

	// TODO: Resolve actions smarter
	return c.resolveActions(actions)
}
