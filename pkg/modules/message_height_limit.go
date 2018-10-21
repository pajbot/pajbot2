// +build csharp

package modules

// To enable the message height limit module, you need .NET Core on your server

// #cgo LDFLAGS: -L../../3rdParty/MessageHeightTwitch/c-interop -lcoreruncommon -ldl -lstdc++
// #include "../../3rdParty/MessageHeightTwitch/c-interop/exports.h"
// #include <stdlib.h>
import "C"

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/utils"
)

var _ pkg.Module = &MessageHeightLimit{}

func floatPtr(v float32) *float32 {
	return &v
}

var messageHeightLimitSpec *moduleSpec

func init() {
	messageHeightLimitSpec = &moduleSpec{
		id:    "message_height_limit",
		name:  "Message height limit",
		maker: NewMessageHeightLimit,

		enabledByDefault: false,

		parameters: map[string]*moduleParameterSpec{
			"HeightLimit": &moduleParameterSpec{
				description:  "Max height of a message before it's timed out",
				defaultValue: floatPtr(95),
			},
		},
	}

	Register(messageHeightLimitSpec)
}

type MessageHeightLimit struct {
	botChannel pkg.BotChannel

	server *server

	HeightLimit floatParameter `json:",omitempty"`
}

func NewMessageHeightLimit() pkg.Module {
	return &MessageHeightLimit{
		server: &_server,

		HeightLimit: floatParameter{
			defaultValue: messageHeightLimitSpec.parameters["HeightLimit"].defaultValue.(*float32),
		},
	}
}

var clrInitialized = false

var messageHeightLimitLibraryInitialized = false
var charMapPath string

func initCLR() error {
	executableDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	fmt.Println("Executable dir", executableDir)
	fmt.Println("os args 0:", os.Args[0])

	clrPath := utils.GetEnv("LIBCOREFOLDER", "/usr/share/dotnet/shared/Microsoft.NETCore.App/2.1.5")

	// Path to our own executable
	clr1 := C.CString(executableDir + "/bot")

	fmt.Println(executableDir)

	// Folder where libcoreclr.so is located
	clr2 := C.CString(clrPath)

	// Path to library we want to use
	clr3 := C.CString(executableDir + "/MessageHeightTwitch.dll")

	var res C.int

	res = C.LoadCLRRuntime(
		clr1,
		clr2,
		clr3)

	C.free(unsafe.Pointer(clr1))
	C.free(unsafe.Pointer(clr2))
	C.free(unsafe.Pointer(clr3))

	if res != 0 {
		return errors.New("Failed to load CLR Runtime")
	}

	charMapPath = executableDir + "/charmap.bin.gz"

	clrInitialized = true

	return nil
}

func initChannel(channelName string) error {
	channel := C.CString(channelName)

	res := C.InitChannel(channel)

	if res != 1 {
		return errors.New("Failed to init Channel " + channelName)
	}

	C.free(unsafe.Pointer(channel))

	return nil
}

func initMessageHeightLimitLibrary() error {
	charMap := C.CString(charMapPath)

	fmt.Println(charMapPath)

	res := C.InitCharMap(charMap)

	C.free(unsafe.Pointer(charMap))

	if res != 1 {
		return errors.New(fmt.Sprintf("Failed to init CharMap: %d", int(res)))
	}

	messageHeightLimitLibraryInitialized = true

	return nil
}

func (m *MessageHeightLimit) Initialize(botChannel pkg.BotChannel, settings []byte) (err error) {
	m.botChannel = botChannel

	if !clrInitialized {
		err = initCLR()
		if err != nil {
			return
		}

		err = initMessageHeightLimitLibrary()
	}

	if err := initChannel(botChannel.ChannelName()); err != nil {
		return err
	}

	if err := loadModule(settings, m); err != nil {
		fmt.Println("Error loading module:", err)
	}

	return
}

func (m *MessageHeightLimit) Disable() error {
	return nil
}

func (m *MessageHeightLimit) Spec() pkg.ModuleSpec {
	return messageHeightLimitSpec
}

func (m *MessageHeightLimit) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m *MessageHeightLimit) OnWhisper(bot pkg.Sender, user pkg.User, message pkg.Message) error {
	return nil
}

func (m *MessageHeightLimit) getHeight(channel pkg.Channel, user pkg.User, message pkg.Message) float32 {
	channelString := C.CString(channel.GetChannel())
	input := C.CString(message.GetText())
	loginName := C.CString(user.GetName())
	displayName := C.CString(user.GetDisplayName())

	var emoteStrings []*C.char

	var te []C.TwitchEmote

	reader := message.GetTwitchReader()
	for reader.Next() {
		emote := reader.Get()
		emoteCode := C.CString(emote.GetName())
		emoteURL := C.CString(fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v1/%s/1.0", emote.GetID()))

		te = append(te, C.TwitchEmote{emoteCode, emoteURL})

		emoteStrings = append(emoteStrings, emoteCode)
		emoteStrings = append(emoteStrings, emoteURL)
	}

	var pArray unsafe.Pointer

	if len(te) > 0 {
		pArray = unsafe.Pointer(&te[0])
	}

	badgeCount := C.int(len(user.GetBadges()))

	height := C.CalculateMessageHeightDirect(
		channelString,
		input,                      // Message text
		loginName,                  // Login name
		displayName,                // Display name
		badgeCount,                 // Badge count
		((*C.TwitchEmote)(pArray)), // Array of emotes
		C.int(len(te)),             // Emote array size
	)

	C.free(unsafe.Pointer(channelString))
	C.free(unsafe.Pointer(input))
	C.free(unsafe.Pointer(loginName))
	C.free(unsafe.Pointer(displayName))

	for _, str := range emoteStrings {
		C.free(unsafe.Pointer(str))
	}

	return float32(height)
}

func (m *MessageHeightLimit) OnMessage(bot pkg.Sender, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	if !messageHeightLimitLibraryInitialized {
		return nil
	}

	if user.GetName() == "gazatu2" {
		return nil
	}

	if user.GetName() == "supibot" {
		return nil
	}

	if user.IsModerator() || user.IsBroadcaster(channel) {
		if strings.HasPrefix(message.GetText(), "!") {
			parts := strings.Split(message.GetText(), " ")
			if parts[0] == "!heightlimit" {
				if len(parts) >= 2 {
					if err := m.HeightLimit.Parse(parts[1]); err != nil {
						bot.Mention(channel, user, err.Error())
						return nil
					}

					bot.Mention(channel, user, "Height limit set to "+utils.Float32ToString(m.HeightLimit.Get()))
					saveModule(m)
				} else {
					bot.Mention(channel, user, "Height limit is "+utils.Float32ToString(m.HeightLimit.Get()))
				}

				return nil
			}

			if parts[0] == "!heighttest" {
				height := m.getHeight(channel, user, message)
				bot.Mention(channel, user, fmt.Sprintf("your message height is %.2f", height))
				return nil
			}
		}
	}

	const maxTimeoutLength = 1800

	height := m.getHeight(channel, user, message)
	// bot.Mention(channel, user, fmt.Sprintf("Message height: %f\n", height))

	if height > m.HeightLimit.Get() {
		timeoutDuration := int(math.Min(math.Pow(float64(height-m.HeightLimit.Get()), 1.2), maxTimeoutLength))
		action.Set(pkg.Timeout{timeoutDuration, fmt.Sprintf("Your message is too tall: %.1f", height)})
	}

	return nil
}
