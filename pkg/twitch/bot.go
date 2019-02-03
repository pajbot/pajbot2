package twitch

import (
	"bufio"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/dankeroni/gotwitch"
	twitch "github.com/gempir/go-twitch-irc"
	"github.com/go-sql-driver/mysql"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/channels"
	"github.com/pajlada/pajbot2/pkg/common"
	"github.com/pajlada/pajbot2/pkg/users"
	"github.com/pajlada/pajbot2/pkg/utils"
	"golang.org/x/oauth2"
)

type ModeState int

var _ pkg.Sender = &Bot{}

const (
	ModeUnset = iota
	ModeEnabled
	ModeDisabled
)

type botFlags struct {
	PermaSubMode bool
}

type BotCredentials struct {
	AccessToken  string
	RefreshToken string
	Expiry       mysql.NullTime
}

// Bot is a wrapper around go-twitch-irc's twitch.Client with a few extra features
type Bot struct {
	*twitch.Client

	app pkg.Application

	TokenSource oauth2.TokenSource

	DatabaseID int

	twitchAccount *User

	QuitChannel chan string

	Flags botFlags

	// Filled in with user IDs when bot is loaded, from SQL
	channelsMutex *sync.Mutex
	channels      []*BotChannel

	// TODO: Store one point server per channel the bot is in. share between bots
	pointServer *PointServer

	ticker *time.Ticker

	userStore   pkg.UserStore
	userContext pkg.UserContext
	streamStore pkg.StreamStore

	pubSub pkg.PubSub

	sql *sql.DB

	IsConnected bool

	onNewChannelJoined func(channelID string)
}

var _ pkg.PubSubConnection = &Bot{}
var _ pkg.PubSubSource = &Bot{}

func NewBot(databaseID int, twitchAccount pkg.TwitchAccount, tokenSource oauth2.TokenSource, app pkg.Application) (*Bot, error) {
	token, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	// TODO(pajlada): share user store between twitch bots
	// TODO(pajlada): mutex lock user store
	b := &Bot{
		app: app,

		TokenSource: tokenSource,
		Client:      twitch.NewClient(twitchAccount.Name(), "oauth:"+token.AccessToken),

		DatabaseID: databaseID,

		twitchAccount: &User{
			name: twitchAccount.Name(),
		},

		channelsMutex: &sync.Mutex{},

		userStore:   app.UserStore(),
		userContext: app.UserContext(),
		streamStore: app.StreamStore(),

		pubSub: app.PubSub(),

		sql: app.SQL(),

		QuitChannel: app.QuitChannel(),
	}

	b.pubSub.Subscribe(b, "Ban")
	b.pubSub.Subscribe(b, "Timeout")
	b.pubSub.Subscribe(b, "Untimeout")

	b.twitchAccount.fillIn(b.userStore)

	return b, nil
}

func (b *Bot) GetTokenSource() oauth2.TokenSource {
	return b.TokenSource
}

func (b *Bot) OnNewChannelJoined(cb func(channelID string)) {
	b.onNewChannelJoined = cb
}

func (b *Bot) IsApplication() bool {
	return true
}

func (b *Bot) Connection() pkg.PubSubConnection {
	return b
}

func (b *Bot) FetchDatabaseID() int {
	return b.DatabaseID
}

func (b *Bot) AuthenticatedUser() pkg.User {
	return nil
}

func (b *Bot) addBotChannel(botChannel *BotChannel) error {
	b.channelsMutex.Lock()
	b.channels = append(b.channels, botChannel)
	b.channelsMutex.Unlock()

	// Initialize bot channel
	return botChannel.Initialize(b)
}

func (b *Bot) getBotChannel(channelID string) (int, *BotChannel) {
	b.channelsMutex.Lock()
	defer b.channelsMutex.Unlock()

	for i, botChannel := range b.channels {
		if botChannel.Channel().GetID() == channelID {
			return i, botChannel
		}
	}

	return -1, nil
}

func (b *Bot) ChannelIDs() (channelIDs []string) {
	b.channelsMutex.Lock()
	defer b.channelsMutex.Unlock()

	for _, botChannel := range b.channels {
		channelIDs = append(channelIDs, botChannel.Channel().GetID())
	}

	return
}

func (b *Bot) InChannel(channelID string) bool {
	b.channelsMutex.Lock()
	defer b.channelsMutex.Unlock()

	for _, botChannel := range b.channels {
		if botChannel.Channel().GetID() == channelID {
			return true
		}
	}

	return false
}

// channelsMutex needs to be locked before calling this function
func (b *Bot) removeBotChannelAtIndex(index int) {
	b.channels = append(b.channels[:index], b.channels[index+1:]...)
}

func (b *Bot) LoadChannels(sql *sql.DB) error {
	const queryF = `SELECT id, twitch_channel_id FROM BotChannel WHERE bot_id=?`

	rows, err := sql.Query(queryF, b.DatabaseID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var channelIDs []string
	var channels []*BotChannel

	for rows.Next() {
		var botChannel BotChannel

		if err = rows.Scan(&botChannel.ID, &botChannel.channel.id); err != nil {
			return err
		}

		channelIDs = append(channelIDs, botChannel.Channel().GetID())
		channels = append(channels, &botChannel)
	}

	m := b.userStore.GetNames(channelIDs)

	for id, name := range m {
		for _, c := range channels {
			if c.Channel().GetID() == id {
				c.channel.SetName(name)
			}
		}
	}

	for _, c := range channels {
		if c.channel.Valid() {
			if err = b.addBotChannel(c); err != nil {
				return nil
			}
		}
	}

	return nil
}

func (b *Bot) Join(channelName string) {
	b.Client.Join(channelName)
	channelID := b.userStore.GetID(channelName)
	if channelID == "" {
		fmt.Println("[pajbot2:pkg/twitch/bot.go] Unable to get ID of channel we just tride to join")
		return
	}

	if b.onNewChannelJoined != nil {
		b.onNewChannelJoined(channelID)
	}
}

func (b *Bot) JoinChannels() {
	b.channelsMutex.Lock()
	defer b.channelsMutex.Unlock()

	for _, c := range b.channels {
		b.Join(c.Channel().GetName())
	}
}

func (b *Bot) GetUserStore() pkg.UserStore {
	return b.userStore
}

func (b *Bot) GetUserContext() pkg.UserContext {
	return b.userContext
}

type emoteReader struct {
	index int

	emotes *[]*common.Emote

	started bool
}

func newEmoteHolder(emotes *[]*common.Emote) *emoteReader {
	return &emoteReader{
		index:  0,
		emotes: emotes,
	}
}

func (h *emoteReader) Next() bool {
	if !h.started {
		h.started = true

		if len(*h.emotes) == 0 {
			return false
		}

		return true
	}

	h.index++

	if h.index >= len(*h.emotes) {
		return false
	}

	return true
}

func (h *emoteReader) Get() pkg.Emote {
	return (*h.emotes)[h.index]
}

// TwitchMessage is a wrapper for twitch.Message with some extra stuff
type TwitchMessage struct {
	twitch.Message

	twitchEmotes      []*common.Emote
	twitchEmoteReader *emoteReader

	bttvEmotes      []*common.Emote
	bttvEmoteReader *emoteReader
	// TODO: BTTV Emotes

	// TODO: FFZ Emotes

	// TODO: Emojis
}

func NewTwitchMessage(message twitch.Message) *TwitchMessage {
	msg := &TwitchMessage{
		Message: message,
	}
	msg.twitchEmoteReader = newEmoteHolder(&msg.twitchEmotes)
	msg.bttvEmoteReader = newEmoteHolder(&msg.bttvEmotes)

	return msg
}

func (m TwitchMessage) GetText() string {
	return m.Text
}

func (m TwitchMessage) GetTwitchReader() pkg.EmoteReader {
	return m.twitchEmoteReader
}

func (m TwitchMessage) GetBTTVReader() pkg.EmoteReader {
	return m.bttvEmoteReader
}

func (m *TwitchMessage) AddBTTVEmote(emote pkg.Emote) {
	m.bttvEmotes = append(m.bttvEmotes, emote.(*common.Emote))
}

// Reply will reply to the message in the same way it received the message
// If the message was received in a twitch channel, reply in that twitch channel.
// IF the message was received in a twitch whisper, reply using twitch whispers.
func (b *Bot) Reply(channel pkg.Channel, user pkg.User, message string) {
	if channel == nil {
		b.Whisper(user, message)
	} else {
		b.Say(channel, message)
	}
}

func (b *Bot) TwitchAccount() pkg.TwitchAccount {
	return b.twitchAccount
}

func (b *Bot) Connected() bool {
	return b.IsConnected
}

func (b *Bot) Say(channel pkg.Channel, message string) {
	b.Client.Say(channel.GetName(), message)
}

func (b *Bot) Mention(channel pkg.Channel, user pkg.User, message string) {
	b.Client.Say(channel.GetName(), "@"+user.GetName()+", "+message)
}

func (b *Bot) Whisper(user pkg.User, message string) {
	b.Client.Whisper(user.GetName(), message)
}

func (b *Bot) Timeout(channel pkg.Channel, user pkg.User, duration int, reason string) {
	if !user.IsModerator() {
		b.Say(channel, fmt.Sprintf(".timeout %s %d %s", user.GetName(), duration, reason))
	}
}

func (b *Bot) Ban(channel pkg.Channel, user pkg.User, reason string) {
	if !user.IsModerator() {
		b.Say(channel, fmt.Sprintf(".ban %s %s", user.GetName(), reason))
	}
}

func (b *Bot) Untimeout(channel pkg.Channel, user pkg.User) {
	if !user.IsModerator() {
		b.Say(channel, fmt.Sprintf(".untimeout %s", user.GetName()))
	}
}

func (b *Bot) HandleWhisper(user twitch.User, rawMessage twitch.Message) {
	twitchUser := users.NewTwitchUser(user, rawMessage.Tags["user-id"])

	// Find out what bot channel this whisper is related to

	parts := strings.Split(rawMessage.Text, " ")
	if len(parts) == 0 {
		return
	}

	channelName := strings.ToLower(utils.FilterChannelName(parts[0]))
	if channelName == "" {
		// No valid channel name was given as context
		// TODO: Pass through to some sort of "global modules"?
		return
	}

	channelID := b.userStore.GetID(channelName)
	if channelID == "" {
		// Context wasn't a valid channel name
		return
	}

	_, botChannel := b.getBotChannel(channelID)
	if botChannel == nil {
		// Whisper was not prefixed with a channel for context, possibly send as a "raw whisper" event?
		return
	}

	rawMessage.Text = strings.Join(parts[1:], " ")

	message := NewTwitchMessage(rawMessage)

	err := botChannel.handleWhisper(twitchUser, message)
	if err != nil {
		fmt.Println("Error occured while forwarding whisper to bot channel:", err)
	}
}

func (b *Bot) HandleMessage(channelName string, user twitch.User, rawMessage twitch.Message) {
	message := NewTwitchMessage(rawMessage)

	twitchUser := users.NewTwitchUser(user, message.Tags["user-id"])

	channel := &channels.TwitchChannel{
		Channel: channelName,
		ID:      rawMessage.Tags["room-id"],
	}

	action := &pkg.TwitchAction{
		Sender:  b,
		Channel: channel,
		User:    twitchUser,
	}

	for _, emote := range rawMessage.Emotes {
		parsedEmote := &common.Emote{
			Name:  emote.Name,
			ID:    emote.ID,
			Count: emote.Count,
			Type:  "twitch",
		}
		message.twitchEmotes = append(message.twitchEmotes, parsedEmote)
	}

	_, botChannel := b.getBotChannel(channel.GetID())
	if botChannel == nil {
		fmt.Println("Message received in channel with id", channel.GetID(), "without having a BotChannel there")
		return
	}

	err := botChannel.handleMessage(twitchUser, message, action)
	if err != nil {
		fmt.Println("Error occured while forwarding message to bot channel:", err)
	}
}

func (b *Bot) HandleRoomstateMessage(channelName string, user twitch.User, rawMessage twitch.Message) {
	if len(rawMessage.Tags) > 2 {
		if channelID, ok := rawMessage.Tags["room-id"]; ok {
			// Joined channel
			b.streamStore.JoinStream(&SimpleAccount{channelID, channelName})
		} else {
			fmt.Println("room-id not set in roomstate message:", rawMessage.Raw)
		}
	}
	subMode := ModeUnset

	channel := &channels.TwitchChannel{
		Channel: channelName,
	}

	if readSubMode, ok := rawMessage.Tags["subs-only"]; ok {
		if readSubMode == "1" {
			subMode = ModeEnabled
		} else {
			subMode = ModeDisabled
		}
	}

	if subMode != ModeUnset {
		if subMode == ModeEnabled {
			// fmt.Println("Submode enabled")
		} else {
			// fmt.Println("Submode disabled")

			if b.Flags.PermaSubMode {
				b.Say(channel, "Perma sub mode is enabled. A mod can type !suboff to disable perma sub mode")
				b.Say(channel, ".subscribers")
			}
		}
	}

	// fmt.Printf("%s - #%s: %#v: %#v\n", b.Name(), channel, user, rawMessage)
}

func (b *Bot) StartChatterPoller() {
	b.ticker = time.NewTicker(5 * time.Minute)
	// defer close ticker lol

	go func() {
		for {
			select {
			case <-b.ticker.C:
				chatters, _, err := gotwitch.GetChattersSimple("pajlada")
				if err != nil {
					continue
				}

				var usernames []string
				usernames = append(usernames, chatters.Moderators...)
				usernames = append(usernames, chatters.Staff...)
				usernames = append(usernames, chatters.Admins...)
				usernames = append(usernames, chatters.GlobalMods...)
				usernames = append(usernames, chatters.Viewers...)
				userIDs := b.GetUserStore().GetIDs(usernames)

				userIDsSlice := make([]string, len(userIDs))
				i := 0
				for _, userID := range userIDs {
					userIDsSlice[i] = userID
					i++
				}

				b.BulkEdit("pajlada", userIDsSlice, 25)
			}
		}
	}()
}

type PointServer struct {
	host string

	conn net.Conn

	ReconnectChannel chan (bool)

	bufferedPayloads [][]byte
}

func (p *PointServer) Write(payload []byte) bool {
	if p.conn != nil {
		_, err := p.conn.Write(payload)
		if err == nil {
			return true
		}

		log.Println("Reconnect????????")
		p.ReconnectChannel <- true
	}

	p.bufferedPayloads = append(p.bufferedPayloads, payload)

	return false
}

func (p *PointServer) Send(command uint8, body []byte) bool {
	bodyLength := make([]byte, 4)
	binary.BigEndian.PutUint32(bodyLength, uint32(len(body)))

	instant := false

	// Write header (Command + Body length)
	instant = p.Write(append([]byte{command}, bodyLength...))

	// Write body
	instant = p.Write(body)

	return instant
}

func (p *PointServer) Read(size int) []byte {
	reader := bufio.NewReader(p.conn)

	response := make([]byte, size)

	_, err := io.ReadFull(reader, response)
	if err != nil {
		p.ReconnectChannel <- true
		return nil
	}

	return response
}

func newPointServer(host string) (*PointServer, error) {
	pointServer := &PointServer{
		host:             host,
		ReconnectChannel: make(chan bool),
	}

	go pointServer.connect()

	return pointServer, nil
}

func (p *PointServer) tryConnect() net.Conn {
	for {
		conn, err := net.Dial("tcp", p.host)

		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		// TODO: Send any buffered messages?

		for _, p := range p.bufferedPayloads {
			conn.Write(p)
		}

		return conn
	}
}

func (p *PointServer) connect() {
	for {
		p.conn = p.tryConnect()

		log.Println("Wait for reconnect channel to proc...")
		<-p.ReconnectChannel
		log.Println("Reconnect... xd")

		p.conn = nil
		log.Println("Reconnect...")
	}

}

func (b *Bot) ConnectToPointServer() (err error) {
	// TODO: read from config file
	b.pointServer, err = newPointServer("localhost:54321")
	if err != nil {
		return
	}

	// TODO: connect once per channel
	b.pointServer.Send(CommandConnect, []byte("pajlada"))

	return
}

const DELIMETER_BYTE = ';'

const (
	CommandConnect   = 0x01
	CommandGetPoints = 0x02
	CommandBulkEdit  = 0x03
	CommandAdd       = 0x04
	CommandRemove    = 0x05
	CommandRank      = 0x06
)

func (b *Bot) GetPoints(channel pkg.Channel, userID string) uint64 {
	if b.pointServer == nil {
		return 0
	}

	bodyPayload := []byte(userID)

	b.pointServer.Send(CommandGetPoints, bodyPayload)
	response := b.pointServer.Read(8)

	if response != nil {
		return binary.BigEndian.Uint64(response)
	}

	return 0
}

func (b *Bot) AddPoints(channel pkg.Channel, userID string, points uint64) (bool, uint64) {
	if b.pointServer == nil {
		return false, 0
	}

	var bodyPayload []byte
	bodyPayload = append(bodyPayload, utils.Uint64ToBytes(points)...)
	bodyPayload = append(bodyPayload, []byte(userID)...)

	b.pointServer.Send(CommandAdd, bodyPayload)

	response := b.pointServer.Read(9)
	userPoints := binary.BigEndian.Uint64(response[1:])

	if response[0] > 0 {
		return false, userPoints
	}

	return true, userPoints
}

func (b *Bot) BulkEdit(channel string, userIDs []string, points int32) {
	if b.pointServer == nil {
		return
	}

	var bodyPayload []byte
	bodyPayload = append(bodyPayload, utils.Int32ToBytes(points)...)
	for _, userID := range userIDs {
		bodyPayload = append(bodyPayload, []byte(userID)...)
		bodyPayload = append(bodyPayload, DELIMETER_BYTE)
	}

	b.pointServer.Send(CommandBulkEdit, bodyPayload)
}

func (b *Bot) RemovePoints(channel pkg.Channel, userID string, points uint64) (bool, uint64) {
	if b.pointServer == nil {
		return false, 0
	}

	var bodyPayload []byte
	bodyPayload = append(bodyPayload, 0x00)
	bodyPayload = append(bodyPayload, utils.Uint64ToBytes(points)...)
	bodyPayload = append(bodyPayload, []byte(userID)...)

	b.pointServer.Send(CommandRemove, bodyPayload)

	response := b.pointServer.Read(9)
	userPoints := binary.BigEndian.Uint64(response[1:])

	if response[0] > 0 {
		return false, userPoints
	}

	return true, userPoints
}

func (b *Bot) ForceRemovePoints(channel pkg.Channel, userID string, points uint64) uint64 {
	if b.pointServer == nil {
		return 0
	}

	var bodyPayload []byte
	bodyPayload = append(bodyPayload, 0x01)
	bodyPayload = append(bodyPayload, utils.Uint64ToBytes(points)...)
	bodyPayload = append(bodyPayload, []byte(userID)...)

	b.pointServer.Send(CommandRemove, bodyPayload)

	response := b.pointServer.Read(9)
	userPoints := binary.BigEndian.Uint64(response[1:])

	return userPoints
}

func (b *Bot) PointRank(channel pkg.Channel, userID string) uint64 {
	if b.pointServer == nil {
		return 0
	}

	var bodyPayload []byte
	bodyPayload = append(bodyPayload, []byte(userID)...)

	b.pointServer.Send(CommandRank, bodyPayload)

	response := b.pointServer.Read(8)
	rank := binary.BigEndian.Uint64(response)

	return rank
}

func FinalMiddleware(bot *Bot, channel pkg.Channel, user pkg.User, message *TwitchMessage, action pkg.Action) {
	// fmt.Printf("Found %d BTTV emotes! %#v", len(message.BTTVEmotes), message.BTTVEmotes)
}

func (b *Bot) MakeUser(username string) pkg.User {
	return users.NewTwitchUser(twitch.User{
		Username:    username,
		DisplayName: username,
	}, b.userStore.GetID(username))
}

func (b *Bot) MakeChannel(channelName string) pkg.Channel {
	return channels.TwitchChannel{
		Channel: channelName,
		ID:      b.userStore.GetID(channelName),
	}
}

func (b *Bot) JoinChannel(channelID string) error {
	const queryF = `INSERT INTO BotChannel (bot_id, twitch_channel_id) VALUES (?, ?)`
	res, err := b.sql.Exec(queryF, b.DatabaseID, channelID)
	if err != nil {
		if common.IsDuplicateKey(err) {
			return errors.New("we have already joined this channel!")
		}

		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	botChannel := &BotChannel{
		ID: id,

		channel: User{
			id: channelID,
		},
	}
	err = botChannel.channel.fillIn(b.userStore)
	if err != nil {
		return err
	}

	b.Join(botChannel.Channel().GetName())

	if err = b.addBotChannel(botChannel); err != nil {
		return err
	}

	if b.onNewChannelJoined != nil {
		b.onNewChannelJoined(channelID)
	}

	return nil
}

func (b *Bot) LeaveChannel(channelID string) error {
	const queryF = `DELETE FROM BotChannel WHERE id=? LIMIT 1;`

	i, botChannel := b.getBotChannel(channelID)
	if botChannel == nil {
		return errors.New("we have not joined this channel")
	}

	b.channelsMutex.Lock()
	defer b.channelsMutex.Unlock()

	b.Depart(botChannel.Channel().GetName())

	res, err := b.sql.Exec(queryF, botChannel.DatabaseID())
	if err != nil {
		return err
	}

	// Delete it from our internal list
	b.removeBotChannelAtIndex(i)

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		// We didn't remove the channel from SQL, s
		return errors.New("unable to remove channel from SQL, but it did exist in our internal storage. sync issue")
	}

	return nil

	return errors.New("we have not joined this channel")
}

func (b *Bot) MessageReceived(source pkg.PubSubSource, topic string, data []byte) error {
	switch topic {
	case "Ban":
		var msg pkg.PubSubBan
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return err
		}
		fmt.Printf("Ban through pubsub: %+v\n", msg)
		b.Ban(b.MakeChannel(msg.Channel), b.MakeUser(msg.Target), msg.Reason)
	case "Timeout":
		var msg pkg.PubSubTimeout
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return err
		}
		fmt.Printf("Timeout through pubsub: %+v\n", msg)
		b.Timeout(b.MakeChannel(msg.Channel), b.MakeUser(msg.Target), int(msg.Duration), msg.Reason)
	case "Untimeout":
		fmt.Printf("untimeout %s\n", string(data))
		var msg pkg.PubSubUntimeout
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return err
		}
		fmt.Printf("Untimeout through pubsub: %+v\n", msg)
		b.Untimeout(b.MakeChannel(msg.Channel), b.MakeUser(msg.Target))
	}
	return nil
}

// Quit quits the entire application
func (b *Bot) Quit(message string) {
	b.channelsMutex.Lock()
	for _, channel := range b.channels {
		channel.Events().Emit("on_quit", nil)
	}
	b.channelsMutex.Unlock()
	time.AfterFunc(250*time.Millisecond, func() { b.QuitChannel <- message })
}
