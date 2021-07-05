package mbase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/eventemitter"
	"github.com/pajbot/pajbot2/pkg/report"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

type Base struct {
	spec pkg.ModuleSpec
	bot  pkg.BotChannel

	ctx    context.Context
	cancel context.CancelFunc

	connections []*eventemitter.Listener

	parameters map[string]pkg.ModuleParameter

	SQL          *sql.DB
	OldSession   *sql.DB
	PubSub       pkg.PubSub
	ReportHolder *report.Holder
}

func New(spec pkg.ModuleSpec, bot pkg.BotChannel, sql, oldSession *sql.DB, pubSub pkg.PubSub, reportHolder *report.Holder) Base {
	b := Base{
		spec: spec,
		bot:  bot,

		parameters: make(map[string]pkg.ModuleParameter),

		SQL:          sql,
		OldSession:   oldSession,
		PubSub:       pubSub,
		ReportHolder: reportHolder,
	}

	parentContext := context.TODO()

	b.ctx, b.cancel = context.WithCancel(parentContext)

	for key, value := range spec.Parameters() {
		b.parameters[key] = value()
	}

	return b
}

func (b *Base) Context() context.Context {
	return b.ctx
}

func (b *Base) BotChannel() pkg.BotChannel {
	return b.bot
}

func (b *Base) MarshalJSON() ([]byte, error) {
	fmt.Println("BASE MARSHAL JSON")

	return nil, nil
}

func (b *Base) Parameters() map[string]pkg.ModuleParameter {
	return b.parameters
}

func (b *Base) LoadSettings(settingsBytes []byte) error {
	if len(b.parameters) == 0 {
		return nil
	}

	if len(settingsBytes) == 0 {
		return nil
	}

	values := map[string]interface{}{}
	err := json.Unmarshal(settingsBytes, &values)
	if err != nil {
		return err
	}

	for key, parameter := range b.parameters {
		if value, ok := values[key]; ok {
			parameter.SetInterface(value)
		}
	}

	return nil
}

func (b *Base) Handle(name string, cb func()) {

}

func (b *Base) Disable() error {
	for _, c := range b.connections {
		c.Disconnected = true
	}

	b.cancel()

	return nil
}

func (b Base) OnWhisper(event pkg.MessageEvent) pkg.Actions {
	return nil
}

func (b Base) OnMessage(event pkg.MessageEvent) pkg.Actions {
	return nil
}

func (b Base) ID() string {
	return b.spec.ID()
}

func (b Base) Type() pkg.ModuleType {
	return b.spec.Type()
}

func (b Base) Priority() int {
	return b.spec.Priority()
}

func (b *Base) SetParameter(key string, value string) error {
	// 1. Find parameter spec (This includes type of the parameter)
	param, ok := b.parameters[key]
	if !ok {
		return fmt.Errorf("no parameter found with the key '%s'", key)
	}

	// 2. Parse `value` according to that parameter spec
	if err := param.Parse(value); err != nil {
		return err
	}

	// 3. If the parameter was updated, update any linked values

	return nil
}

// SetParameterResponse is a helper function which calls SetParameter and returns actions in a standardized format
func (b *Base) SetParameterResponse(key, value string, event pkg.MessageEvent) pkg.Actions {
	if err := b.SetParameter(key, value); err != nil {
		return twitchactions.Mention(event.User, err.Error())
	}

	err := b.Save()
	if err != nil {
		return twitchactions.Mentionf(event.User, "an error occurred while saving parameters for module %s, key %s: %s", b.ID(), key, err.Error())
	}
	return twitchactions.Mentionf(event.User, "%s set to %s", key, value)
}

func (b *Base) save() error {
	parameters := map[string]interface{}{}
	for key, param := range b.Parameters() {
		if !param.HasValue() {
			continue
		}

		if !param.HasBeenSet() {
			continue
		}

		parameters[key] = param.Get()
	}

	bytes, err := json.Marshal(parameters)
	if err != nil {
		return err
	}

	const queryF = `
INSERT INTO
	bot_channel_module
	(bot_channel_id, module_id, settings)
	VALUES ($1, $2, $3)
ON CONFLICT (bot_channel_id, module_id) DO UPDATE SET settings=$3`

	_, err = b.SQL.Exec(queryF, b.bot.DatabaseID(), b.ID(), bytes) // GOOD
	if err != nil {
		return err
	}

	return nil
}

func (b *Base) Save() error {
	err := b.save()
	if err != nil {
		log.Printf("Error saving module %s: %s\n", b.ID(), err.Error())
		return err
	}

	return nil
}

func (b *Base) Listen(event string, cb interface{}, prio int) error {
	conn, err := b.bot.Events().Listen(event, cb, prio)
	if err != nil {
		return err
	}

	b.connections = append(b.connections, conn)

	return nil
}
