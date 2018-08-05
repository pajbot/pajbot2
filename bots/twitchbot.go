package bots

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/dankeroni/gotwitch"
	twitch "github.com/gempir/go-twitch-irc"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/channels"
	pb2twitch "github.com/pajlada/pajbot2/pkg/twitch"
	"github.com/pajlada/pajbot2/pkg/users"
	"github.com/pajlada/pajbot2/pkg/utils"
	"github.com/pajlada/pajbot2/redismanager"
)

type ModeState int

var _ pkg.Sender = &TwitchBot{}

const (
	ModeUnset = iota
	ModeEnabled
	ModeDisabled
)

type botFlags struct {
	PermaSubMode bool
}

// TwitchBot is a wrapper around go-twitch-irc's twitch.Client with a few extra features
type TwitchBot struct {
	*twitch.Client

	Name    string
	handler Handler

	QuitChannel chan string

	Flags botFlags

	Redis *redismanager.RedisManager

	Modules []pkg.Module

	// TODO: Store one point server per channel the bot is in. share between bots
	pointServer *PointServer

	ticker *time.Ticker

	userStore *pb2twitch.UserStore
}

func NewTwitchBot(client *twitch.Client) *TwitchBot {
	// TODO(pajlada): share user store between twitch bots
	// TODO(pajlada): mutex lock user store
	return &TwitchBot{
		Client:    client,
		userStore: pb2twitch.NewUserStore(),
	}
}

func (b *TwitchBot) GetUserStore() pkg.UserStore {
	return b.userStore
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
func (b *TwitchBot) Reply(channel pkg.Channel, user pkg.User, message string) {
	if channel == nil {
		b.Whisper(user, message)
	} else {
		b.Say(channel, message)
	}
}

func (b *TwitchBot) Say(channel pkg.Channel, message string) {
	b.Client.Say(channel.GetChannel(), message)
}

func (b *TwitchBot) Mention(channel pkg.Channel, user pkg.User, message string) {
	b.Client.Say(channel.GetChannel(), "@"+user.GetName()+", "+message)
}

func (b *TwitchBot) Whisper(user pkg.User, message string) {
	b.Client.Whisper(user.GetName(), message)
}

func (b *TwitchBot) Timeout(channel pkg.Channel, user pkg.User, duration int, reason string) {
	// Empty string in UserType means a normal user
	if !user.IsModerator() {
		b.Say(channel, fmt.Sprintf(".timeout %s %d %s", user.GetName(), duration, reason))
	}
}

// SetHandler sets the handler to message at the bottom of the list
func (b *TwitchBot) SetHandler(handler Handler) {
	b.handler = handler
}

func (b *TwitchBot) HandleWhisper(user twitch.User, rawMessage twitch.Message) {
	message := NewTwitchMessage(rawMessage)

	twitchUser := users.NewTwitchUser(user, message.Tags["user-id"])

	action := &pkg.TwitchAction{
		Sender: b,
		User:   twitchUser,
	}

	b.handler.HandleMessage(b, nil, twitchUser, message, action)

	if pkg.VerboseMessages {
		log.Printf("%s - @%s(%s): %s", b.Name, twitchUser.DisplayName, twitchUser.Username, message.Text)
	}
}

func (b *TwitchBot) HandleMessage(channelName string, user twitch.User, rawMessage twitch.Message) {
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

	b.handler.HandleMessage(b, channel, twitchUser, message, action)

	if pkg.VerboseMessages {
		log.Printf("%s - #%s: %s(%s): %s", b.Name, channel, twitchUser.DisplayName, twitchUser.Username, message.Text)
	}
}

func (b *TwitchBot) HandleRoomstateMessage(channelName string, user twitch.User, rawMessage twitch.Message) {
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
			log.Printf("Submode enabled")
		} else {
			log.Printf("Submode disabled")

			if b.Flags.PermaSubMode {
				b.Say(channel, "Perma sub mode is enabled. A mod can type !suboff to disable perma sub mode")
				b.Say(channel, ".subscribers")
			}
		}
	}

	log.Printf("%s - #%s: %#v: %#v", b.Name, channel, user, rawMessage)
}

// Quit quits the entire application
func (b *TwitchBot) Quit(message string) {
	b.QuitChannel <- message
}
func onHTTPError(statusCode int, statusMessage, errorMessage string) {
	log.Println("HTTPERROR: ", errorMessage)
}

func onInternalError(err error) {
	log.Printf("internal error: %s", err)
}

func (b *TwitchBot) StartChatterPoller() {
	b.ticker = time.NewTicker(15 * time.Second)
	// defer close ticker lol

	go func() {
		for {
			select {
			case <-b.ticker.C:
				onSuccess := func(chatters gotwitch.Chatters) {
					usernames := []string{}
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

					b.BulkEdit("pajlada", userIDsSlice, 5)
				}

				gotwitch.GetChatters("pajlada", onSuccess, onHTTPError, onInternalError)
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
		log.Println("Trying to connect to", p.host)
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

func (b *TwitchBot) ConnectToPointServer() (err error) {
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
)

func (b *TwitchBot) GetPoints(channel pkg.Channel, user pkg.User) uint64 {
	bodyPayload := []byte(user.GetID())

	b.pointServer.Send(CommandGetPoints, bodyPayload)
	response := b.pointServer.Read(8)

	if response != nil {
		return binary.BigEndian.Uint64(response)
	}

	return 0
}

func (b *TwitchBot) AddPoints(channel pkg.Channel, userID string, points uint64) (bool, uint64) {
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

func (b *TwitchBot) BulkEdit(channel string, userIDs []string, points int32) {
	var bodyPayload []byte
	bodyPayload = append(bodyPayload, utils.Int32ToBytes(points)...)
	for _, userID := range userIDs {
		bodyPayload = append(bodyPayload, []byte(userID)...)
		bodyPayload = append(bodyPayload, DELIMETER_BYTE)
	}

	b.pointServer.Send(CommandBulkEdit, bodyPayload)
}

func (b *TwitchBot) RemovePoints(channel pkg.Channel, userID string, points uint64) (bool, uint64) {
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

func (b *TwitchBot) ForceRemovePoints(channel pkg.Channel, userID string, points uint64) uint64 {
	var bodyPayload []byte
	bodyPayload = append(bodyPayload, 0x01)
	bodyPayload = append(bodyPayload, utils.Uint64ToBytes(points)...)
	bodyPayload = append(bodyPayload, []byte(userID)...)

	b.pointServer.Send(CommandRemove, bodyPayload)

	response := b.pointServer.Read(9)
	userPoints := binary.BigEndian.Uint64(response[1:])

	return userPoints
}

func (b *TwitchBot) AddModule(module pkg.Module) {
	if module == nil {
		return
	}

	if err := module.Register(); err != nil {
		log.Printf("Error registering module(%s): %s\n", module.Name(), err.Error())
		return
	}

	b.Modules = append(b.Modules, module)
}
