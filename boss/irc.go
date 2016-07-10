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

	"github.com/dghubble/go-twitter/twitter"
	"github.com/pajlada/pajbot2/parser"
	"github.com/pajlada/pajbot2/pbtwitter"
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
	brokerHost  string
	brokerPass  string
	brokerLogin string
	pass        string
	nick        string
	conn        net.Conn
	join        chan string
	ReadChan    chan string
	SendChan    chan string
	bots        map[string]chan common.Msg
	redis       *redismanager.RedisManager
	sql         *sqlmanager.SQLManager
	twitter     *pbtwitter.Client
	quit        chan string
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
	if irc.brokerLogin != "" {
		irc.SendRaw(conn, "LOGIN "+irc.brokerLogin)
		irc.SendRaw(conn, "PASS "+irc.pass)
	} else if irc.pass != "" {
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
				m := parser.Parse(line)
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
					if b := irc.getBot(m.Channel); b != nil {
						b <- m
					} else {
						log.Debugf("default No channel for message (chan: %s)", m.Channel)
					}
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
	irc.twitter.Bots[channel] = &pbtwitter.Bot{
		Following: []string{"nuulss", "pajlada", "pajbot"},
		Stream:    make(chan *twitter.Tweet, 5),
		Client:    irc.twitter,
	}
	read := make(chan common.Msg, 10)
	newbot := bot.Config{
		Quit:     irc.quit,
		Channel:  channel,
		ReadChan: read,
		SendChan: irc.SendChan,
		Join:     irc.join,
		Redis:    irc.redis,
		SQL:      irc.sql,
		Twitter:  irc.twitter.Bots[channel],
	}
	irc.bots[channel] = read
	b := bot.NewBot(newbot)
	_modules := []bot.Module{
		&modules.Banphrase{},
		&modules.Command{},
		&modules.Pyramid{},
		&modules.Quit{},
		&modules.SubAnnounce{},
		&modules.MyInfo{},
		&modules.Test{},
		&modules.Points{},
		&modules.Top{},
		&modules.Raffle{},
	}
	b.Modules = _modules
	for _, mod := range b.Modules {
		mod.Init(b)
	}
	go b.Init()
}

/*
JoinChannel joins a twitch chat and creates a new bot if there isnt already one
*/
func (irc *Irc) JoinChannel(channel string) {
	channel = strings.ToLower(channel)
	irc.Lock()
	defer irc.Unlock()
	if _, ok := irc.bots[channel]; !ok {
		irc.NewBot(channel)
		irc.SendRaw(irc.conn, "JOIN #"+channel)
	}
}

/*
PartChannel leaves a twitch channel
but the bot is still running and able to post in that chat
TODO: proper Bot.Close() that stops all its go routines
*/
func (irc *Irc) PartChannel(channel string) {
	channel = strings.ToLower(channel)
	irc.Lock()
	defer irc.Unlock()
	if bot, ok := irc.bots[channel]; ok {
		delete(irc.bots, channel)
		irc.SendRaw(irc.conn, "PART #"+channel)
		close(bot)
		log.Debug("CLOSED BOT IN", channel)

	}
}

/*
JoinChannels joins a list of channels, given as a string slice
can also be used to part channels when using the prefix "PART "
this might be confusing but a part channel is overkill imo
*/
func (irc *Irc) JoinChannels() {
	for line := range irc.join {
		log.Debug(line)
		if strings.HasPrefix(line, "PART ") {
			channel := strings.Split(line, " ")[1]
			irc.PartChannel(channel)
			log.Debug("PART CHANNEL", channel)
		} else {
			irc.JoinChannel(line)
		}
	}
}

/*
Init initalizes shit.

TODO: This should just create the Irc object. You should have to call
irc.Run() manually I think. or irc.Start()?
*/
func Init(config *common.Config) *Irc {
	irc := &Irc{
		brokerHost:  *config.BrokerHost,
		brokerPass:  *config.BrokerPass,
		brokerLogin: config.BrokerLogin,
		pass:        config.Pass,
		nick:        config.Nick,
		ReadChan:    make(chan string, 10),
		SendChan:    make(chan string, 10),
		join:        make(chan string, 5),
		bots:        make(map[string]chan common.Msg),
		redis:       redismanager.Init(config),
		sql:         sqlmanager.Init(config),
		quit:        config.Quit,
	}
	irc.twitter = pbtwitter.Init(config, irc.redis)
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
	go irc.JoinChannels()
	for _, channel := range config.Channels {
		irc.join <- channel
	}
	return irc
}

// XXX: rename to getBotChannel? idk
func (irc *Irc) getBot(channel string) chan common.Msg {
	if b, ok := irc.bots[channel]; ok {
		return b
	}

	return nil
}
