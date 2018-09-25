package twitch

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/dankeroni/gotwitch"
	twitch "github.com/gempir/go-twitch-irc"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/channels"
	"github.com/pajlada/pajbot2/pkg/common"
	"github.com/pajlada/pajbot2/pkg/pubsub"
	"github.com/pajlada/pajbot2/pkg/users"
	"github.com/pajlada/pajbot2/pkg/utils"
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

// Bot is a wrapper around go-twitch-irc's twitch.Client with a few extra features
type Bot struct {
	*twitch.Client

	Name    string
	handler Handler

	QuitChannel chan string

	Flags botFlags

	Modules []pkg.Module

	// TODO: Store one point server per channel the bot is in. share between bots
	pointServer *PointServer

	ticker *time.Ticker

	userStore pkg.UserStore

	pubSub *pubsub.PubSub
}

var _ pubsub.Connection = &Bot{}

func NewBot(client *twitch.Client, pubSub *pubsub.PubSub, userStore pkg.UserStore) *Bot {
	// TODO(pajlada): share user store between twitch bots
	// TODO(pajlada): mutex lock user store
	b := &Bot{
		Client:    client,
		userStore: userStore,

		pubSub: pubSub,
	}

	pubSub.Subscribe(b, "Ban")
	pubSub.Subscribe(b, "Untimeout")

	return b
}

func (b *Bot) GetUserStore() pkg.UserStore {
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
func (b *Bot) Reply(channel pkg.Channel, user pkg.User, message string) {
	if channel == nil {
		b.Whisper(user, message)
	} else {
		b.Say(channel, message)
	}
}

func (b *Bot) Say(channel pkg.Channel, message string) {
	b.Client.Say(channel.GetChannel(), message)
}

func (b *Bot) Mention(channel pkg.Channel, user pkg.User, message string) {
	b.Client.Say(channel.GetChannel(), "@"+user.GetName()+", "+message)
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

// SetHandler sets the handler to message at the bottom of the list
func (b *Bot) SetHandler(handler Handler) {
	b.handler = handler
}

func (b *Bot) HandleWhisper(user twitch.User, rawMessage twitch.Message) {
	message := NewTwitchMessage(rawMessage)

	twitchUser := users.NewTwitchUser(user, message.Tags["user-id"])

	action := &pkg.TwitchAction{
		Sender: b,
		User:   twitchUser,
	}

	b.handler.HandleMessage(b, nil, twitchUser, message, action)

	if pkg.VerboseMessages {
		fmt.Printf("%s - @%s(%s): %s", b.Name, twitchUser.DisplayName, twitchUser.Username, message.Text)
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

	b.handler.HandleMessage(b, channel, twitchUser, message, action)

	if pkg.VerboseMessages {
		fmt.Printf("%s - #%s: %s(%s): %s", b.Name, channel, twitchUser.DisplayName, twitchUser.Username, message.Text)
	}
}

func (b *Bot) HandleRoomstateMessage(channelName string, user twitch.User, rawMessage twitch.Message) {
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
			fmt.Println("Submode enabled")
		} else {
			fmt.Println("Submode disabled")

			if b.Flags.PermaSubMode {
				b.Say(channel, "Perma sub mode is enabled. A mod can type !suboff to disable perma sub mode")
				b.Say(channel, ".subscribers")
			}
		}
	}

	// fmt.Printf("%s - #%s: %#v: %#v\n", b.Name, channel, user, rawMessage)
}

// Quit quits the entire application
func (b *Bot) Quit(message string) {
	b.QuitChannel <- message
}
func onHTTPError(statusCode int, statusMessage, errorMessage string) {
	log.Println("HTTPERROR: ", errorMessage)
}

func onInternalError(err error) {
	fmt.Printf("internal error: %s", err)
}

func (b *Bot) StartChatterPoller() {
	b.ticker = time.NewTicker(5 * time.Minute)
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

					b.BulkEdit("pajlada", userIDsSlice, 25)
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

func (b *Bot) AddModule(module pkg.Module) {
	if module == nil {
		return
	}

	if err := module.Register(); err != nil {
		fmt.Printf("Error registering module(%s): %s\n", module.Name(), err.Error())
		return
	}

	b.Modules = append(b.Modules, module)
}

// TODO: Code under here should be generalized, or moved into their own module
func HandleCommands(next Handler) Handler {
	return HandlerFunc(func(bot *Bot, channel pkg.Channel, user pkg.User, message *TwitchMessage, action pkg.Action) {
		if user.IsModerator() || user.IsBroadcaster(channel) || user.GetName() == "pajlada" || user.GetName() == "karl_kons" || user.GetName() == "fourtf" {
			if strings.HasPrefix(message.Text, "!xd") {
				bot.Reply(channel, user, "XDDDDDDDDDD")
				return
			}

			if strings.HasPrefix(message.Text, "!myuserid") {
				bot.Say(channel, fmt.Sprintf("@%s, your user ID is %s", user.GetName(), user.GetID()))
				return
			}

			if strings.HasPrefix(message.Text, "!whisperme") {
				fmt.Printf("Send whisper!")
				bot.Say(channel, "@"+user.GetName()+", I just sent you a whisper with the text \"hehe\" :D")
				bot.Whisper(user, "hehe")
				return
			}

			if strings.HasPrefix(message.Text, "!modme") {
				bot.Say(channel, ".mod "+user.GetName())
				bot.Say(channel, "Modded")
				return
			}

			if strings.HasPrefix(message.Text, "!unmodme") {
				bot.Say(channel, ".unmod "+user.GetName())
				bot.Say(channel, "Unmodded")
				return
			}

			if strings.HasPrefix(message.Text, "!pb2quit") {
				bot.Reply(channel, user, "Quitting...")
				time.AfterFunc(time.Millisecond*500, func() {
					bot.Quit("Quit because pajlada said so")
				})
				return
			}

			if strings.HasPrefix(message.Text, "!emoteonly") {
				bot.Say(channel, ".emoteonly")
				return
			}

			if strings.HasPrefix(message.Text, "!emoteonlyoff") || message.Text == "TriHard TriHard TriHard forsenE pajaCool TriHard" {
				bot.Say(channel, ".emoteonlyoff")
				return
			}

			if strings.HasPrefix(message.Text, "!subon") {
				if bot.Flags.PermaSubMode {
					bot.Say(channel, "Permanent subscribers mode is already enabled")
					return
				}

				bot.Flags.PermaSubMode = true

				bot.Say(channel, ".subscribers")
				bot.Say(channel, "Permanent subscribers mode has been enabled")
				return
			}

			if strings.HasPrefix(message.Text, "!suboff") {
				if !bot.Flags.PermaSubMode {
					bot.Say(channel, "Permanent subscribers mode is not enabled")
					return
				}

				bot.Flags.PermaSubMode = false

				bot.Say(channel, ".subscribersoff")
				bot.Say(channel, "Permanent subscribers mode has been disabled")
				return
			}
		}

		next.HandleMessage(bot, channel, user, message, action)
	})
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

func (b *Bot) MakeChannel(channel string) pkg.Channel {
	return channels.TwitchChannel{
		Channel: channel,
		ID:      b.userStore.GetID(channel),
	}
}

func (b *Bot) MessageReceived(topic string, data []byte) error {
	switch topic {
	case "Ban":
		var msg pkg.PubSubBan
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return err
		}
		b.Ban(b.MakeChannel(msg.Channel), b.MakeUser(msg.Target), msg.Reason)
	case "Untimeout":
		fmt.Printf("untimeout %s\n", string(data))
		var msg pkg.PubSubUntimeout
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return err
		}
		b.Untimeout(b.MakeChannel(msg.Channel), b.MakeUser(msg.Target))
	}
	return nil
}
