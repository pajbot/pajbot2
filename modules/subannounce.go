package modules

import (
	"encoding/json"
	"time"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/common/basemodule"
)

/*
SubAnnounce xD
*/
type SubAnnounce struct {
	basemodule.BaseModule
	NewSubMessage string `json:"new_sub_message"`
	ResubMessage  string `json:"resub_message"`
	NewSubWhisper string `json:"new_sub_whisper"`
	ResubWhisper  string `json:"new_sub_whisper"`
	WhisperDelay  int    `json:"whisper_delay"`
}

// Ensure the module implements the interface properly
var _ Module = (*SubAnnounce)(nil)

func (module *SubAnnounce) parseSettings(jsonData []byte) {
	type Alias SubAnnounce
	/*
		aux := struct {
			*Alias
		}{
			Alias: (*Alias)(module),
		}
	*/
	if err := json.Unmarshal(jsonData, module); err != nil {
		log.Error(err)
	}
}

// Init xD
func (module *SubAnnounce) Init(bot *bot.Bot) (string, bool) {
	module.SetDefaults("sub-announce")
	module.EnabledDefault = true
	module.ParseState(bot.Redis, bot.Channel.Name)

	module.parseSettings(module.FetchSettings(bot.Redis, bot.Channel.Name))

	return "sub-announce", true
}

// DeInit xD
func (module *SubAnnounce) DeInit(b *bot.Bot) {

}

// Check xD
func (module *SubAnnounce) Check(b *bot.Bot, m *common.Msg, action *bot.Action) error {
	if m.Type == common.MsgSub {
		data := map[string]string{
			"username": m.User.Name,
		}

		if module.NewSubMessage != "" {
			action.Response = common.FormatString(module.NewSubMessage, data)
			action.Stop = true
		}
		if module.NewSubWhisper != "" {
			delay := time.Second * time.Duration(module.WhisperDelay)
			time.AfterFunc(delay, func() {
				whisperMessage := common.FormatString(module.NewSubWhisper, data)
				b.Whisper(m.User.Name, whisperMessage)
			})
		}
	} else if m.Type == common.MsgReSub {
		data := map[string]string{
			"username": m.User.Name,
			"months":   m.Tags["msg-param-months"],
		}

		if module.ResubMessage != "" {
			action.Response = common.FormatString(module.ResubMessage, data)
			action.Stop = true
		}
		if module.ResubWhisper != "" {
			delay := time.Second * time.Duration(module.WhisperDelay)
			time.AfterFunc(delay, func() {
				whisperMessage := common.FormatString(module.ResubWhisper, data)
				b.Whisper(m.User.Name, whisperMessage)
			})
		}
	}
	return nil
}
