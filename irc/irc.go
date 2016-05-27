package irc

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"time"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/helper"
)

/*
The Irc object contains all data xD
*/
type Irc struct {
	server   string
	port     string
	pass     string
	nick     string
	readConn map[net.Conn][]string
	sendConn map[net.Conn][]int
	ReadChan chan string
	SendChan chan string
	channels map[string]net.Conn
	bots     map[string]chan bot.Msg
}

/*
SendRaw sends a raw message to the given connection.
The only thing it appends is \r\n
*/
func (irc *Irc) SendRaw(s net.Conn, line string) {
	fmt.Fprint(s, line+"\r\n")
}

func (irc *Irc) newConn(send bool) {
	conn, err := net.Dial("tcp", irc.server+":"+irc.port)
	if err != nil {
		fmt.Println("Error connecting to the IRC servers:", err)
		return
	}
	if irc.pass != "" {
		irc.SendRaw(conn, "PASS "+irc.pass)
	}
	irc.SendRaw(conn, "NICK "+irc.nick)
	irc.SendRaw(conn, "CAP REQ twitch.tv/tags")
	if send {
		irc.sendConn[conn] = make([]int, 30)
		go irc.keepAlive(conn)
	} else {
		irc.readConn[conn] = make([]string, 0)
		go irc.readConnection(conn)
	}
	fmt.Println("connected")
}

func (irc *Irc) getSendConn() net.Conn {
	var conn net.Conn
	for c := range irc.sendConn {
		if helper.Sum(irc.sendConn[c]) < 15 {
			conn = c
			break
		}
	}
	if conn == nil {
		irc.newConn(true)
		conn = irc.getSendConn()
	}
	return conn
}

func (irc *Irc) send() {
	for {
		msg := <-irc.SendChan
		conn := irc.getSendConn()
		irc.SendRaw(conn, msg)
		fmt.Println("sent: " + msg)
		irc.sendConn[conn][29]++
	}
}

func (irc *Irc) rateLimit() {
	for {
		for conn, s := range irc.sendConn {
			newS := append(s[1:], 0)
			irc.sendConn[conn] = newS
		}
		time.Sleep(1 * time.Second)
	}
}

func (irc *Irc) keepAlive(conn net.Conn) {
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	for {
		line, err := tp.ReadLine()
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(line, "PING") {
			irc.SendRaw(conn, strings.Replace(line, "PING", "PONG", 1))
		}
	}

}
func (irc *Irc) readConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	for {
		line, err := tp.ReadLine()
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(line, "PING") {
			irc.SendRaw(conn, strings.Replace(line, "PING", "PONG", 1))
		} else if strings.Contains(line, "PRIVMSG") || strings.Contains(line, "WHISPER") {
			m := Parse(line)
			irc.bots[m.Channel] <- m
		}
	}
}

/*
JoinChannel creates a new bot in the given channel

TODO: Only create a new bot if one doesn't already exist in the channel
*/
func (irc *Irc) JoinChannel(channel string) {
	conn := irc.getReadconn()
	irc.SendRaw(conn, "JOIN #"+channel)
	irc.readConn[conn] = append(irc.readConn[conn], channel)
	read := make(chan bot.Msg)
	newbot := &bot.Config{
		Channel:  channel,
		ReadChan: read,
		SendChan: irc.SendChan,
	}
	irc.bots[channel] = read
	go bot.NewBot(*newbot)

}

/*
JoinChannels joins a list of channels, given as a string slice
*/
func (irc *Irc) JoinChannels(channels []string) {
	for i := range channels {
		irc.JoinChannel(channels[i])
		time.Sleep(300 * time.Millisecond)
	}
}

func (irc *Irc) getReadconn() net.Conn {
	var conn net.Conn
	for c, channels := range irc.readConn {
		if len(channels) < 50 {
			conn = c
			break
		}
	}
	if conn == nil {
		irc.newConn(false)
		conn = irc.getReadconn()
	}
	return conn
}

/*
Init initalizes shit.

TODO: This should just create the Irc object. You should have to call
irc.Run() manually I think. or irc.Start()?
*/
func Init(pass string, nick string) Irc {
	irc := &Irc{
		server:   "irc.chat.twitch.tv",
		port:     "80",
		pass:     pass,
		nick:     nick,
		readConn: make(map[net.Conn][]string),
		sendConn: make(map[net.Conn][]int),
		ReadChan: make(chan string, 10),
		SendChan: make(chan string, 10),
		bots:     make(map[string]chan bot.Msg),
	}
	irc.newConn(true)
	irc.newConn(false)
	go irc.send()
	go irc.rateLimit()
	return *irc
}
