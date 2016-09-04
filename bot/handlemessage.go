package bot

import "github.com/pajlada/pajbot2/common"

/*
Handle attempts to handle the given message
*/
func (b *Bot) Handle(msg common.Msg) {
	b.parseBttvEmotes(&msg)
	common.ParseEmojis(&msg)
	oldUser := msg.User
	defer b.Redis.UpdateUser(b.Channel.Name, &msg.User, &oldUser)
	action := &Action{}
	log.Debugf("%s # %s :%s", msg.Channel, msg.User.DisplayName, msg.Text)
	for _, module := range b.EnabledModules {
		// If user level is above module bypass level
		//   then don't call Check here
		// TODO: implement above shit
		module.Check(b, &msg, action)

		if action.Response != "" {
			b.Say(action.Response)
			action.Response = "" // delete Response
		}

		if action.Stop {
			return
		}
	}

	if b.Channel.Online {
		msg.User.OnlineMessageCount++
	} else {
		msg.User.OfflineMessageCount++
	}
	msg.User.TotalMessageCount++
}
