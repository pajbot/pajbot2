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
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"github.com/pajbot/pajbot2/pkg"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
	"github.com/pajbot/utils"
)

func init() {
	Register("message_height_limit", func() pkg.ModuleSpec {
		return &Spec{
			id:    "message_height_limit",
			name:  "Message height limit",
			maker: NewMessageHeightLimit,

			moduleType: pkg.ModuleTypeFilter,

			enabledByDefault: false,

			parameters: map[string]pkg.ModuleParameterSpec{
				"HeightLimit": func() pkg.ModuleParameter {
					return newFloatParameter(parameterSpec{
						Description:  "Max height of a message before it's timed out",
						DefaultValue: float32(95.0),
					})
				},
				"AsciiArtOnly": func() pkg.ModuleParameter {
					return newBoolParameter(parameterSpec{
						Description:  "Only attempt to catch ascii art",
						DefaultValue: false,
					})
				},
			},
		}
	})
}

var _ pkg.Module = &MessageHeightLimit{}

type MessageHeightLimit struct {
	mbase.Base

	HeightLimit float32

	AsciiArtOnly bool

	userViolationCount map[string]int
}

func NewMessageHeightLimit(b mbase.Base) pkg.Module {
	m := &MessageHeightLimit{
		Base: b,

		userViolationCount: make(map[string]int),
	}

	m.Parameters()["HeightLimit"].Link(&m.HeightLimit)
	m.Parameters()["AsciiArtOnly"].Link(&m.AsciiArtOnly)

	// FIXME
	m.Initialize()

	return m
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

	fmt.Println("init channel", channelName)
	res := C.InitChannel(channel)
	fmt.Println("done")

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

func (m *MessageHeightLimit) Initialize() {
	var err error
	fmt.Println("Initializing message height limit")

	if !clrInitialized {
		fmt.Println("init clr..")
		err = initCLR()
		if err != nil {
			return
		}
		fmt.Println("done")

		err = initMessageHeightLimitLibrary()
		fmt.Println("done init height limit library")
	}

	fmt.Println("init channel")
	if err := initChannel(m.BotChannel().ChannelName()); err != nil {
		log.Println("Error initializing channel:", err)
		return
	}
	fmt.Println("done")
}

func (m *MessageHeightLimit) getHeight(channel pkg.Channel, user pkg.User, message pkg.Message) float32 {
	channelString := C.CString(channel.GetName())
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

func (m *MessageHeightLimit) OnMessage(event pkg.MessageEvent) pkg.Actions {
	if !messageHeightLimitLibraryInitialized {
		return nil
	}

	user := event.User
	message := event.Message

	if user.HasPermission(m.BotChannel().Channel(), pkg.PermissionImmuneToMessageLimits) {
		return nil
	}

	if user.IsModerator() || user.HasPermission(m.BotChannel().Channel(), pkg.PermissionModeration) {
		if strings.HasPrefix(message.GetText(), "!") {
			parts := strings.Split(message.GetText(), " ")
			if parts[0] == "!heightlimit" {
				if len(parts) >= 2 {
					if err := m.SetParameter("HeightLimit", parts[1]); err != nil {
						return twitchactions.Mention(user, err.Error())
					}

					m.Save()
					return twitchactions.Mentionf(user, "Height limit set to %.0f", m.HeightLimit)
				}

				return twitchactions.Mentionf(user, "Height limit is %.0f", m.HeightLimit)
			}

			if parts[0] == "!heighttest" {
				height := m.getHeight(m.BotChannel().Channel(), user, message)
				return twitchactions.Mentionf(user, "your message height is %.2f", height)
			}

			if parts[0] == "!heightlimitonasciionly" {
				if len(parts) >= 2 {
					if err := m.SetParameter("AsciiArtOnly", parts[1]); err != nil {
						return twitchactions.Mention(user, err.Error())
					}

					m.Save()
					return twitchactions.Mentionf(user, "Height limit module set to act on ascii art only: %v", m.AsciiArtOnly)
				}

				return twitchactions.Mentionf(user, "Height limit module is set to act on ascii art only: %v", m.AsciiArtOnly)
			}
		}
	}

	const minTimeoutLength = 10
	const maxTimeoutLength = 1800

	height := m.getHeight(m.BotChannel().Channel(), user, message)

	if height <= m.HeightLimit {
		return nil
	}

	actions := &twitchactions.Actions{}

	// Message height is too tall
	messageLength := len([]rune(message.GetText()))
	var fitsIn7Bit int
	var doesntFitIn7Bit int
	for _, r := range message.GetText() {
		if r > 0x7a || r < 0x20 {
			doesntFitIn7Bit++
		} else {
			fitsIn7Bit++
		}
	}

	fmt.Printf("Message length: %d. Fits: %d. Don't fit: %d\n", messageLength, doesntFitIn7Bit, fitsIn7Bit)
	var ratio float32
	ratio = float32(doesntFitIn7Bit) / float32(messageLength)
	var reason string
	userViolations := 0
	timeoutDuration := int(math.Min(math.Pow(float64(height-m.HeightLimit), 1.2), maxTimeoutLength))
	if ratio > 0.5 {
		timeoutDuration = timeoutDuration + 90
	} else {
		if m.AsciiArtOnly {
			// Do not deal with tall non-ascii-art messages
			return nil
		}
	}

	timeoutDuration = utils.MaxInt(minTimeoutLength, timeoutDuration)

	const reasonFmt = `Your message is too tall: %.1f - %.3f (%d)`

	if ratio > 0.5 && height > 140.0 {
		m.userViolationCount[user.GetID()] = m.userViolationCount[user.GetID()] + 1
		userViolations = m.userViolationCount[user.GetID()]
		timeoutDuration = timeoutDuration * userViolations
		timeoutDuration = utils.MinInt(3600*24*7, timeoutDuration)
		actions.Whisper(user, fmt.Sprintf("Your message is too long and contains too many non-ascii characters. Your next timeout will be multiplied by %d", userViolations))
	}

	reason = fmt.Sprintf(reasonFmt, height, ratio, userViolations)
	actions.Timeout(user, time.Duration(timeoutDuration)*time.Second).SetReason(reason)

	return actions
}
