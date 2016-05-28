package bot

/*
Module xD
*/
type Module interface {
	Check(bot *Bot, msg *Msg, action *Action) error
}
