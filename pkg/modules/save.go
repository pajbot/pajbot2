package modules

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pajbot/pajbot2/pkg"
)

func saveModule(module pkg.Module) error {
	if module == nil {
		return errors.New("saveModule: module may not be nil")
	}

	botChannel := module.BotChannel()
	if botChannel == nil {
		return errors.New("saveModule: No bot channel specified for module " + module.ID())
	}

	parameters := map[string]interface{}{}
	for key, param := range module.Parameters() {
		if !param.HasValue() {
			continue
		}

		if !param.HasBeenSet() {
			continue
		}

		parameters[key] = param.Get()
	}

	b, err := json.Marshal(parameters)
	if err != nil {
		return err
	}

	const queryF = `
INSERT INTO
	BotChannelModule
	(bot_channel_id, module_id, settings)
	VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE settings=?`

	_, err = _server.sql.Exec(queryF, botChannel.DatabaseID(), module.ID(), b, b)
	if err != nil {
		return err
	}

	return nil
}

func loadParameters(settings []byte) (map[string]interface{}, error) {
	parameters := map[string]interface{}{}
	err := json.Unmarshal(settings, &parameters)
	if err != nil {
		return nil, err
	}

	return parameters, nil
}

func loadFloat(p map[string]interface{}, key string, value interface{}) {
	pValue, ok := p[key]
	if !ok {
		return
	}

	fmt.Println("Attempt to load", pValue, "into", key)
}

func loadModule(settings []byte, module interface{}) error {
	if module == nil {
		return errors.New("loadModule: module may not be nil")
	}

	if len(settings) == 0 {
		return nil
	}

	fmt.Println("Loading", string(settings))

	return json.Unmarshal(settings, module)
}
