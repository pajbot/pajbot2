package irc

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"nuulsbot/src/bot"
	"strings"
	"time"
)

type Irc struct {
	server   string
	port     string
	pass     string
	nick     string
	readconn map[net.Conn][]string
	sendconn map[net.Conn][]int
	Readchan chan string
	Sendchan chan string
	channels map[string]net.Conn
	bots     map[string]chan bot.Msg
}

func (irc *Irc) SendRaw(s net.Conn, line string) {
	fmt.Fprint(s, line+"\r\n")
}

func (irc *Irc) newConn(send bool) {
	conn, _ := net.Dial("tcp", irc.server+":"+irc.port)
	if irc.pass != "" {
		irc.SendRaw(conn, "PASS "+irc.pass)
	}
	irc.SendRaw(conn, "NICK "+irc.nick)
	irc.SendRaw(conn, "CAP REQ twitch.tv/tags")
	if send {
		irc.sendconn[conn] = make([]int, 30)
		go irc.keepAlive(conn)
	} else {
		irc.readconn[conn] = make([]string, 0)
		go irc.readConnection(conn)
	}
	fmt.Println("connected")
}

func (irc *Irc) getSendConn() net.Conn {
	var conn net.Conn
	for c := range irc.sendconn {
		if Sum(irc.sendconn[c]) < 15 {
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
		msg := <-irc.Sendchan
		conn := irc.getSendConn()
		irc.SendRaw(conn, msg)
		fmt.Println("sent: " + msg)
		irc.sendconn[conn][29]++
	}
}

func (irc *Irc) rateLimit() {
	for {
		for conn, s := range irc.sendconn {
			newS := append(s[1:], 0)
			irc.sendconn[conn] = newS
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

func (irc *Irc) JoinChannel(channel string) {
	conn := irc.getReadconn()
	irc.SendRaw(conn, "JOIN #"+channel)
	irc.readconn[conn] = append(irc.readconn[conn], channel)
	read := make(chan bot.Msg)
	newbot := &bot.BotConfig{
		Channel:  channel,
		Readchan: read,
		Sendchan: irc.Sendchan,
	}
	irc.bots[channel] = read
	go bot.NewBot(*newbot)

}

func (irc *Irc) JoinChannels(channels []string) {
	for i := range channels {
		irc.JoinChannel(channels[i])
		time.Sleep(300 * time.Millisecond)
	}
}

func (irc *Irc) getReadconn() net.Conn {
	var conn net.Conn
	for c, channels := range irc.readconn {
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

func Init(pass string, nick string) Irc {
	irc := &Irc{
		server:   "irc.chat.twitch.tv",
		port:     "80",
		pass:     pass,
		nick:     nick,
		readconn: make(map[net.Conn][]string),
		sendconn: make(map[net.Conn][]int),
		Readchan: make(chan string, 10),
		Sendchan: make(chan string, 10),
		bots:     make(map[string]chan bot.Msg),
	}
	irc.newConn(true)
	irc.newConn(false)
	go irc.send()
	go irc.rateLimit()
	return *irc
}
