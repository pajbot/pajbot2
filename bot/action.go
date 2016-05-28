package bot

/*
An Action is the value that is returned from every module
and every method that is called on each irc message
*/
type Action struct {
	Stop     bool
	Response string
}
