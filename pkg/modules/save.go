package modules

import (
	"encoding/json"
	"errors"

	"github.com/pajbot/pajbot2/pkg"
)

func saveModule(module pkg.Module) error {
	if module == nil {
		return errors.New("saveModule: module may not be nil")
	}

	botChannel := module.BotChannel()
	if botChannel == nil {
		return errors.New("saveModule: No bot channel specified for module " + module.Spec().ID())
	}

	b, err := json.Marshal(module)
	if err != nil {
		return err
	}

	const queryF = `
INSERT INTO
	BotChannelModule
	(bot_channel_id, module_id, settings)
	VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE settings=?`

	_, err = _server.sql.Exec(queryF, botChannel.DatabaseID(), module.Spec().ID(), b, b)
	if err != nil {
		return err
	}

	return nil
}

func loadModule(settings []byte, module pkg.Module) error {
	if module == nil {
		return errors.New("loadModule: module may not be nil")
	}

	return json.Unmarshal(settings, module)
}
