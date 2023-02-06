package twitch

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/nicklaw5/helix/v2"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/eventemitter"
	"github.com/pajbot/pajbot2/pkg/modules"
	"github.com/pajbot/utils"
)

var _ pkg.BotChannel = &BotChannel{}

// BotChannel describes a bot account that's connected to a channel
type BotChannel struct {
	stream pkg.Stream

	ID int64

	channel User
	BotUser User

	initialized bool

	// Enabled modules
	modules      []pkg.Module
	modulesMutex sync.Mutex

	sql *sql.DB

	// The eventEmitter is a bot-channel specific event emitter where things related to the channel, e.g. the "stream store" can push data that modules can read
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

func (c *BotChannel) Untimeout(user pkg.User) {
	c.bot.Untimeout(&c.channel, user)
}

func (c *BotChannel) Ban(user pkg.User, reason string) {
	c.bot.Ban(&c.channel, user, reason)
}

func (c *BotChannel) Unban(user pkg.User) {
	c.bot.Unban(&c.channel, user)
}

func (c *BotChannel) DeleteMessage(message string) {
	c.bot.DeleteMessage(&c.channel, message)
}

func (c *BotChannel) DatabaseID() int64 {
	return c.ID
}

func (c *BotChannel) Events() *eventemitter.EventEmitter {
	return c.eventEmitter
}

func (c *BotChannel) GetName() string {
	return c.channel.Name()
}

func (c *BotChannel) GetID() string {
	return c.channel.ID()
}

func (c *BotChannel) ChannelID() string {
	return c.channel.ID()
}

func (c *BotChannel) ChannelName() string {
	return c.channel.Name()
}

func (c *BotChannel) Stream() pkg.Stream {
	return c.stream
}

// We assume that modulesMutex is locked already
func (c *BotChannel) sortModules() {
	sort.Slice(c.modules, func(i, j int) bool {
		a := c.modules[i]
		b := c.modules[j]

		if a.Priority() == b.Priority() {
			return a.Type() > b.Type()
		}

		return a.Priority() < b.Priority()
	})
}

func (c *BotChannel) getSettingsForModule(moduleID string) ([]byte, error) {
	const queryF = `
SELECT
	settings
FROM
	bot_channel_module
WHERE
	bot_channel_id=$1 AND module_id=$2`

	row := c.sql.QueryRow(queryF, c.DatabaseID(), moduleID) // GOOD

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
	bot_channel_module
	(bot_channel_id, module_id, enabled)
	VALUES ($1, $2, $3)
ON CONFLICT (bot_channel_id, module_id) DO UPDATE SET enabled=$3`

	_, err := c.sql.Exec(queryF, c.DatabaseID(), moduleID, state) // GOOD
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// EnableModule enables a module with the given id
// Returns an error if the module is not registered
// Returns an error if the module is already enabled
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

// DisableModule disables a module with the given id
// Returns an error if the module is already disabled
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

func (c *BotChannel) GetModule(moduleID string) (pkg.Module, error) {
	c.modulesMutex.Lock()
	defer c.modulesMutex.Unlock()

	for _, m := range c.modules {
		if m.ID() == moduleID {
			return m, nil
		}
	}

	return nil, errors.New("no module with this ID found")
}

func (c *BotChannel) Initialize(b *Bot) error {
	if c.initialized {
		return errors.New("bot channel is already initialized")
	}

	c.bot = b
	c.sql = b.sql

	c.initialized = true

	c.eventEmitter = eventemitter.New()

	c.loadModules()

	c.eventEmitter.Emit("on_join", nil)

	c.stream = b.streamStore.GetStream(&User{
		id: c.channel.id,
	})

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
	const queryF = `SELECT id, module_id, enabled, settings FROM bot_channel_module WHERE bot_channel_id=$1`

	rows, err := c.sql.Query(queryF, c.DatabaseID()) // GOOD
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

	availableModules := modules.List()

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

func (c *BotChannel) OnModules(cb func(module pkg.Module) pkg.Actions, stop bool) (actions []pkg.Actions) {
	c.modulesMutex.Lock()
	defer c.modulesMutex.Unlock()

	for _, module := range c.modules {
		// TODO: This could potentially be run in steps now. maybe all modules with same priority are run together
		if moduleActions := cb(module); moduleActions != nil {
			actions = append(actions, moduleActions)
			if stop && moduleActions.StopPropagation() {
				return
			}
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

		for _, unmute := range action.Unmutes() {
			switch unmute.Type() {
			case pkg.MuteTypeTemporary:
				c.Untimeout(unmute.User())
			case pkg.MuteTypePermanent:
				c.Unban(unmute.User())
			}
		}

		for _, delete := range action.Deletes() {
			c.DeleteMessage(delete.Message())
		}

		for _, message := range action.Messages() {
			c.Say(message.Evaluate())
		}

		for _, whisper := range action.Whispers() {
			c.bot.Whisper(whisper.User(), whisper.Content())
		}
	}

	return nil
}

func (c *BotChannel) HandleMessage(user pkg.User, message pkg.Message) error {
	c.eventEmitter.Emit("on_msg", nil)

	event := pkg.MessageEvent{
		BaseEvent: pkg.BaseEvent{
			UserStore: c.bot.GetUserStore(),
		},
		User:    user,
		Message: message,
		Channel: c,
	}

	actions := c.OnModules(func(module pkg.Module) pkg.Actions {
		return module.OnMessage(event)
	}, true)

	return c.resolveActions(actions)
}

func (c *BotChannel) HandleEventSubNotification(notification pkg.TwitchEventSubNotification) error {
	// We might want to handle some events here, so read the subscription type and handle before forwarding to modules

	switch notification.Subscription.Type {
	case helix.EventSubTypeStreamOnline:
		var onlineEvent helix.EventSubStreamOnlineEvent
		err := json.NewDecoder(bytes.NewReader(notification.Event)).Decode(&onlineEvent)
		if err != nil {
			return nil
		}

		// Fake the helix stream as much as we can
		s := &helix.Stream{
			ID:        onlineEvent.ID,
			Type:      onlineEvent.Type,
			StartedAt: onlineEvent.StartedAt.Time,
		}
		c.stream.Update(s)
	case helix.EventSubTypeStreamOffline:
		c.stream.Update(nil)
	}

	event := pkg.EventSubNotificationEvent{
		BaseEvent: pkg.BaseEvent{
			UserStore: c.bot.GetUserStore(),
		},
		Notification: notification,
	}

	actions := c.OnModules(func(module pkg.Module) pkg.Actions {
		return module.OnEventSubNotification(event)
	}, true)

	return c.resolveActions(actions)
}

func (c *BotChannel) handleWhisper(user pkg.User, message *TwitchMessage) error {
	event := pkg.MessageEvent{
		BaseEvent: pkg.BaseEvent{
			UserStore: c.bot.GetUserStore(),
		},
		User:    user,
		Message: message,
		Channel: c,
	}

	actions := c.OnModules(func(module pkg.Module) pkg.Actions {
		return module.OnWhisper(event)
	}, true)

	return c.resolveActions(actions)
}

func (c *BotChannel) SetSubscribers(state bool) error {
	c.bot.UpdateChatSettings(&c.channel, &helix.UpdateChatSettingsParams{
		SubscriberMode: &state,
	})
	return nil
}

func (c *BotChannel) SetUniqueChat(state bool) error {
	c.bot.UpdateChatSettings(&c.channel, &helix.UpdateChatSettingsParams{
		UniqueChatMode: &state,
	})
	return nil
}

func (c *BotChannel) SetEmoteOnly(state bool) error {
	c.bot.UpdateChatSettings(&c.channel, &helix.UpdateChatSettingsParams{
		EmoteMode: &state,
	})
	return nil
}

func (c *BotChannel) SetSlowMode(state bool, durationS int) error {
	c.bot.UpdateChatSettings(&c.channel, &helix.UpdateChatSettingsParams{
		SlowMode:         &state,
		SlowModeWaitTime: &durationS,
	})
	return nil
}

func (c *BotChannel) SetFollowerMode(state bool, durationM int) error {
	c.bot.UpdateChatSettings(&c.channel, &helix.UpdateChatSettingsParams{
		FollowerMode:         &state,
		FollowerModeDuration: &durationM,
	})
	return nil
}

func (c *BotChannel) SetNonModChatDelay(state bool, durationS int) error {
	c.bot.UpdateChatSettings(&c.channel, &helix.UpdateChatSettingsParams{
		NonModeratorChatDelay:         &state,
		NonModeratorChatDelayDuration: &durationS,
	})
	return nil
}
