package main

import (
	"time"

	"github.com/nuuls/pajbot2/config"
	"github.com/nuuls/pajbot2/irc"
)

func main() {
	config := config.GetConfig()
	irc := irc.Init(config.Pass, config.Nick)
	irc.JoinChannel("nuuls")
	//irc.JoinChannel("forsenlol")
	irc.JoinChannel("pajlada")
	// irc.JoinChannel("lirik")
	// irc.JoinChannel("timthetatman")
	// irc.JoinChannel("witwix")
	// irc.JoinChannel("watchmeblink1")
	for {
		time.Sleep(5 * time.Second)
	}
}
