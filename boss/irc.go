package boss

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
	"sync"
	"time"

	"github.com/pajlada/pajbot2/redismanager"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/helper"
	"github.com/pajlada/pajbot2/modules"
)

/*
The Irc object contains all data xD
*/
type Irc struct {
	sync.Mutex
	server   string
	port     string
	pass     string
	nick     string
	readConn map[net.Conn][]string
	sendConn map[net.Conn][]int
	ReadChan chan string
	SendChan chan string
	channels map[string]net.Conn
	bots     map[string]chan common.Msg
	redis    *redismanager.RedisManager
	parser   *Parse
	quit     chan string
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
	irc.Lock()
	defer irc.Unlock()
	// wait for connection, this should be done better but we're gonna use
	// relaybroker anyways so it should be fine for now
	time.Sleep(500 * time.Millisecond)
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
		irc.Lock()
		irc.sendConn[conn][29]++
		irc.Unlock()
	}
}

func (irc *Irc) rateLimit() {
	for {
		for conn, s := range irc.sendConn {
			newS := append(s[1:], 0)
			irc.Lock()
			irc.sendConn[conn] = newS
			irc.Unlock()
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
			log.Println("connection died", err)
			delete(irc.sendConn, conn)
			return
		}
		if strings.HasPrefix(line, "PING") {
			irc.SendRaw(conn, strings.Replace(line, "PING", "PONG", 1))
		}
	}
}

func (irc *Irc) GetGlobalUser(m *common.Msg) {
	u := &common.GlobalUser{}
	irc.redis.GetGlobalUser(m.Channel, &m.User, u)
	if m.Type == common.MsgWhisper {
		m.Channel = u.Channel
	}
}

func (irc *Irc) readConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	for {
		line, err := tp.ReadLine()
		if err != nil {
			log.Println("connection died", err)
			irc.newConn(false)
			irc.JoinChannels(irc.readConn[conn])
			delete(irc.readConn, conn)
			return
		}
		if strings.HasPrefix(line, "PING") {
			irc.SendRaw(conn, strings.Replace(line, "PING", "PONG", 1))
		} else if strings.Contains(line, "PRIVMSG") || strings.Contains(line, "WHISPER") {
			m := irc.parser.Parse(line)
			// throw away its own and other useless msgs
			if m.Type != common.MsgThrowAway && m.User.Name != irc.nick {
				irc.GetGlobalUser(&m)
				irc.bots[m.Channel] <- m
			}
		}
	}
}

// NewBot creates a new bot in the given channel
func (irc *Irc) NewBot(channel string) {
	read := make(chan common.Msg)
	newbot := bot.Config{
		Quit:     irc.quit,
		Channel:  channel,
		ReadChan: read,
		SendChan: irc.SendChan,
		Redis:    irc.redis,
	}
	irc.bots[channel] = read
	commandModule := &modules.Command{}
	// TODO: This should be generalized (and optional if possible)
	// Could that be based on module type?
	// If module.@type == 'NeedsInit' { (cast)module.Init() }
	commandModule.Init()
	_modules := []bot.Module{
		&modules.Banphrase{},
		commandModule,
		&modules.Pyramid{},
		&modules.Quit{},
	}
	b := bot.NewBot(newbot, _modules)
	go b.Init()
}

/*
JoinChannel joins a twitch chat and creates a new bot if there isnt already one
*/
func (irc *Irc) JoinChannel(channel string) {
	conn := irc.getReadconn()
	irc.SendRaw(conn, "JOIN #"+channel)
	irc.Lock()
	defer irc.Unlock()
	if _, ok := irc.bots[channel]; !ok {
		irc.readConn[conn] = append(irc.readConn[conn], channel)
		irc.NewBot(channel)
	}

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
func Init(config *common.Config) *Irc {
	irc := &Irc{
		server:   "irc.chat.twitch.tv",
		port:     "80",
		pass:     config.Pass,
		nick:     config.Nick,
		readConn: make(map[net.Conn][]string),
		sendConn: make(map[net.Conn][]int),
		ReadChan: make(chan string, 10),
		SendChan: make(chan string, 10),
		bots:     make(map[string]chan common.Msg),
		redis:    redismanager.Init(config),
		parser:   &Parse{},
		quit:     config.Quit,
	}
	irc.newConn(true)
	irc.newConn(false)
	go irc.send()
	go irc.rateLimit()
	go irc.JoinChannels(config.Channels)
	return irc
}
