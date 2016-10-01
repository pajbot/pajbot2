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
)

var log = plog.GetLogger()

// IRCConfig xD
type IRCConfig struct {
	BrokerHost  string
	BrokerPass  string
	BrokerLogin string
	Redis       *redismanager.RedisManager
	SQL         *sqlmanager.SQLManager
	Twitter     *pbtwitter.Client
	Quit        chan string
	Silent      bool
}

/*
The Irc object contains all data xD
*/
type Irc struct {
	sync.Mutex
	BotAccountID int
	brokerHost   string
	brokerPass   string
	brokerLogin  string
	pass         string
	nick         string
	conn         net.Conn
	join         chan string
	ReadChan     chan string
	SendChan     chan string
	Bots         map[string]*bot.Bot
	Redis        *redismanager.RedisManager
	SQL          *sqlmanager.SQLManager
	twitter      *pbtwitter.Client
	quit         chan string
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
		// This is used in silent Bots
		<-irc.SendChan
	}
}

// GetGlobalUser fills in the global user in the message from Redis
func (irc *Irc) GetGlobalUser(m *common.Msg) {
	u := &common.GlobalUser{}
	irc.Redis.GetGlobalUser(m.Channel, &m.User, u)
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
						b.Read <- m
					} else {
						log.Debugf("No channel for message (chan: %s)", m.Channel)
					}
				case common.MsgSub:
					// Post sub to sub channel
					log.Debugf("%s just subbed!", m.User.DisplayName)
					if b := irc.getBot(m.Channel); b != nil {
						b.Read <- m
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
						b.Read <- m
					} else {
						log.Debugf("MsgReSub No channel for message (chan: %s)", m.Channel)
					}
				case common.MsgThrowAway:
					// Do nothing
					break
				default:
					if b := irc.getBot(m.Channel); b != nil {
						b.Read <- m
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
		BotAccountID: irc.BotAccountID,
		Quit:         irc.quit,
		Channel:      channel,
		ReadChan:     read,
		SendChan:     irc.SendChan,
		Join:         irc.join,
		Redis:        irc.Redis,
		SQL:          irc.SQL,
		Twitter:      irc.twitter.Bots[channel],
	}
	b := bot.NewBot(newbot)
	irc.Bots[channel] = b

	// Populate bot.AllModules with an instance of all available modules
	modulesInit(b)

	// Call Init on all modules, then push in all enabled modules into
	// the bot.EnabledModules list
	modulesLoad(b)

	go b.Init()
}

/*
JoinChannel joins a twitch chat and creates a new bot if there isnt already one
*/
func (irc *Irc) JoinChannel(channel string) {
	channel = strings.ToLower(channel)
	irc.Lock()
	defer irc.Unlock()
	if _, ok := irc.Bots[channel]; !ok {
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
	if bot, ok := irc.Bots[channel]; ok {
		delete(irc.Bots, channel)
		irc.SendRaw(irc.conn, "PART #"+channel)
		close(bot.Read)
		log.Debug("CLOSED BOT IN", channel)

	}
}

/*
JoinChannels listens to requests from the `irc.join` channel and
joins those channels.
If the message starts with "PART" we instead leave that channel.
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
InitIRCConnection initalizes shit.

TODO: This should just create the Irc object. You should have to call
irc.Run() manually I think. or irc.Start()?
*/
func InitIRCConnection(config IRCConfig, botAccount common.DBUser) *Irc {
	irc := &Irc{
		BotAccountID: botAccount.ID,
		brokerHost:   config.BrokerHost,
		brokerPass:   config.BrokerPass,
		brokerLogin:  config.BrokerLogin,
		pass:         "oauth:" + botAccount.TwitchCredentials.AccessToken,
		nick:         botAccount.Name,
		ReadChan:     make(chan string, 10),
		SendChan:     make(chan string, 10),
		join:         make(chan string, 5),
		Bots:         make(map[string]*bot.Bot),
		Redis:        config.Redis,
		SQL:          config.SQL,
		twitter:      config.Twitter,
		quit:         config.Quit,
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

	// Start a goroutine which handles joining and parting from channels
	go irc.JoinChannels()

	channels, err := common.FetchAllChannels(irc.SQL, botAccount.ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Got %d channels", len(channels))

	hasOwnChannel := false
	for _, channel := range channels {
		log.Debug(channel.Name)
		if channel.Name == botAccount.Name {
			hasOwnChannel = true
		}

		if !channel.Enabled {
			log.Debugf("Skipping %s cuz not enabled", channel.Name)
			continue
		}

		irc.join <- channel.Name
	}

	if !hasOwnChannel {
		// Create our own channel, then use InsertNewToSQL
		ownChannel := &common.Channel{
			Name:  botAccount.Name,
			BotID: botAccount.ID,
		}
		ownChannel.InsertNewToSQL(irc.SQL)

		irc.join <- ownChannel.Name
	}

	return irc
}

// XXX: rename to getBotChannel? idk
func (irc *Irc) getBot(channel string) *bot.Bot {
	if b, ok := irc.Bots[channel]; ok {
		return b
	}

	return nil
}
