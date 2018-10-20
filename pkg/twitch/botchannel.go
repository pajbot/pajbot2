package twitch

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/modules"
)

var _ pkg.BotChannel = &BotChannel{}

type BotChannel struct {
	ID int64

	Channel User
	BotUser User

	initialized bool

	// Enabled modules
	modules []pkg.Module

	sql *sql.DB
}

func (c *BotChannel) DatabaseID() int64 {
	return c.ID
}

func (c *BotChannel) ChannelID() string {
	return c.Channel.ID
}

func (c *BotChannel) ChannelName() string {
	return c.Channel.Name
}

func (c *BotChannel) Initialize(b *Bot) error {
	if c.initialized {
		return errors.New("bot channel is already initialized")
	}

	c.sql = b.sql

	c.initialized = true

	c.loadModules()

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

	for _, cfg := range moduleConfigs {
		fmt.Printf("cfg: %+v\n", cfg)
	}

	availableModules := modules.Modules()
	fmt.Println("Available modules:", availableModules)
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
			fmt.Println("Enabling module", spec.Name())

			module := spec.Maker()()
			err := module.Initialize(c, settings)
			if err != nil {
				fmt.Printf("Error loading module '%s': %s\n", spec.ID(), err.Error())
				continue
			}

			c.modules = append(c.modules, module)
		}
		// Fetch config for this module from SQL
	}
}

func (c *BotChannel) forwardToModules(bot pkg.Sender, channel pkg.Channel, user pkg.User, message *TwitchMessage, action pkg.Action) error {

	for _, module := range c.modules {
		var err error
		if channel == nil {
			err = module.OnWhisper(bot, user, message)
		} else {
			err = module.OnMessage(bot, channel, user, message, action)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
