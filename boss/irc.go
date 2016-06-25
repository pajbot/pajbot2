package boss

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"sync"

	"github.com/pajlada/pajbot2/plog"
	"github.com/pajlada/pajbot2/redismanager"
	"github.com/pajlada/pajbot2/sqlmanager"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/modules"
)

var log = plog.GetLogger()

/*
The Irc object contains all data xD
*/
type Irc struct {
	sync.Mutex
	brokerHost string
	brokerPass string
	pass       string
	nick       string
	conn       net.Conn
	ReadChan   chan string
	SendChan   chan string
	bots       map[string]chan common.Msg
	redis      *redismanager.RedisManager
	sql        *sqlmanager.SQLManager
	parser     *parse
	quit       chan string
}

/*
SendRaw sends a raw message to the given connection.
The only thing it appends is \r\n
*/
func (irc *Irc) SendRaw(s net.Conn, line string) {
	fmt.Fprint(s, line+"\r\n")
}

func (irc *Irc) newConn() error {
	if irc.conn != nil {
		// A connection already exists
		return nil
	}
	conn, err := net.Dial("tcp", irc.brokerHost)
	if err != nil {
		return errors.New("Error connecting to the IRC servers:" + err.Error())
	}
	if irc.pass != "" {
		irc.SendRaw(conn, "PASS "+irc.brokerPass+";"+irc.pass)
	}
	irc.SendRaw(conn, "NICK "+irc.nick)
	go irc.readConnection(conn)
	irc.conn = conn
	log.Debug("connected")
	return nil
}

func (irc *Irc) send() {
	for {
		msg := <-irc.SendChan
		irc.SendRaw(irc.conn, msg)
		log.Debugf("sent: %s", msg)
	}
}

func (irc *Irc) dontSend() {
	for {
		// Do nothing with the messages.
		// This is used in silent bots
		<-irc.SendChan
	}
}

// GetGlobalUser fills in the global user in the message from redis
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
	readChan := make(chan string)
	running := true
	go func() {
		var line string
		for running {
			line = <-readChan
			if strings.HasPrefix(line, "PING") {
				irc.SendRaw(conn, strings.Replace(line, "PING", "PONG", 1))
			} else {
				m := irc.parser.Parse(line)
				// throw away its own and other useless msgs
				if m.User.Name == irc.nick {
					// Throw away its own messages
					continue
				}
				switch m.Type {
				case common.MsgPrivmsg, common.MsgWhisper:
					irc.GetGlobalUser(&m)
					if b := irc.getBot(m.Channel); b != nil {
						b <- m
					} else {
						log.Debugf("No channel for message (chan: %s)", m.Channel)
					}
				case common.MsgSub:
					// Post sub to sub channel
					log.Debugf("%s just subbed!", m.User.DisplayName)
					if b := irc.getBot(m.Channel); b != nil {
						b <- m
					} else {
						log.Debugf("MsgSub No channel for message (chan: %s)", m.Channel)
					}
				case common.MsgReSub:
					months, err := strconv.Atoi(m.Tags["msg-param-months"])
					if err != nil {
						// ERROR
						break
					}
					// TODO: check room-id
					log.Debugf("%s just resubbed for %d months", m.User.DisplayName, months)

					if b := irc.getBot(m.Channel); b != nil {
						b <- m
					} else {
						log.Debugf("MsgReSub No channel for message (chan: %s)", m.Channel)
					}
				case common.MsgThrowAway:
					// Do nothing
					break
				default:
					log.Debugf("Unhandled message[%d]: %s\n", m.Type, m.Message)
				}
			}
		}
	}()
	defer func() {
		running = false
		close(readChan)
	}()
	for {
		line, err := tp.ReadLine()
		if err != nil {
			log.Debug("connection died", err)
			irc.newConn()
			//irc.JoinChannels(irc.readConn[conn])
			return
		}
		readChan <- line
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
	commandModule.Init(irc.sql)
	banphraseModule := &modules.Banphrase{}
	banphraseModule.Init(irc.sql)
	_modules := []bot.Module{
		banphraseModule,
		commandModule,
		&modules.Pyramid{},
		&modules.Quit{},
		&modules.SubAnnounce{},
		&modules.MyInfo{},
	}
	b := bot.NewBot(newbot, _modules)
	go b.Init()
}

/*
JoinChannel joins a twitch chat and creates a new bot if there isnt already one
*/
func (irc *Irc) JoinChannel(channel string) {
	irc.Lock()
	defer irc.Unlock()
	if _, ok := irc.bots[channel]; !ok {
		irc.NewBot(channel)
		irc.SendRaw(irc.conn, "JOIN #"+channel)
	}
}

/*
JoinChannels joins a list of channels, given as a string slice
*/
func (irc *Irc) JoinChannels(channels []string) {
	for _, channel := range channels {
		irc.JoinChannel(channel)
	}
}

/*
Init initalizes shit.

TODO: This should just create the Irc object. You should have to call
irc.Run() manually I think. or irc.Start()?
*/
func Init(config *common.Config) *Irc {
	irc := &Irc{
		brokerHost: *config.BrokerHost,
		brokerPass: *config.BrokerPass,
		pass:       config.Pass,
		nick:       config.Nick,
		ReadChan:   make(chan string, 10),
		SendChan:   make(chan string, 10),
		bots:       make(map[string]chan common.Msg),
		redis:      redismanager.Init(config),
		sql:        sqlmanager.Init(config),
		parser:     &parse{},
		quit:       config.Quit,
	}
	err := irc.newConn()
	if err != nil {
		// Right now we just fatally exit the bot
		// You're personally responsible for restarting the bot if it crashes
		log.Fatal(err)
	}
	if config.Silent {
		go irc.dontSend()
	} else {
		go irc.send()
	}
	go irc.JoinChannels(config.Channels)
	return irc
}

// XXX: rename to getBotChannel? idk
func (irc *Irc) getBot(channel string) chan common.Msg {
	if b, ok := irc.bots[channel]; ok {
		return b
	}

	return nil
}
