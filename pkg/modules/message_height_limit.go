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
	"os"
	"path/filepath"
	"unsafe"

	"github.com/pajlada/pajbot2/pkg"
)

var _ pkg.Module = &MessageHeightLimit{}

type MessageHeightLimit struct {
	server *server
}

func NewMessageHeightLimit() *MessageHeightLimit {
	return &MessageHeightLimit{
		server: &_server,
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

	clrLibraryFolder, clrLibraryFolderSet := os.LookupEnv("LIBCOREFOLDER")
	if !clrLibraryFolderSet {
		clrLibraryFolder = "/opt/dotnet/shared/Microsoft.NETCore.App/2.1.5/"
	}

	// Path to our own executable
	clr1 := C.CString(executableDir + "/bot")

	// Folder where libcoreclr.so is located
	clr2 := C.CString(clrLibraryFolder)

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

func initMessageHeightLimitLibrary() error {
	charMap := C.CString(charMapPath)
	channel := C.CString("forsen")

	fmt.Println(charMapPath)

	res := C.InitMessageHeightTwitch(charMap, channel)

	C.free(unsafe.Pointer(charMap))
	C.free(unsafe.Pointer(channel))

	if res != 0 {
		return errors.New("Failed to init MessageHeightTwitch")
	}

	messageHeightLimitLibraryInitialized = true

	return nil
}

func (m *MessageHeightLimit) Register() (err error) {
	if !clrInitialized {
		err = initCLR()
		if err != nil {
			return
		}

		err = initMessageHeightLimitLibrary()
	}

	return
}

func (m *MessageHeightLimit) Name() string {
	return "MessageHeightLimit"
}

func (m *MessageHeightLimit) OnWhisper(bot pkg.Sender, user pkg.User, message pkg.Message) error {
	return nil
}

func (m *MessageHeightLimit) getHeight(user pkg.User, message pkg.Message) float32 {
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
		input,                      // Message text
		loginName,                  // Login name
		displayName,                // Display name
		badgeCount,                 // Badge count
		((*C.TwitchEmote)(pArray)), // Array of emotes
		C.int(len(te)),             // Emote array size
	)

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

	if channel.GetChannel() != "forsen" {
		return nil
	}

	const heightLimit = 160

	height := m.getHeight(user, message)
	bot.Mention(channel, user, fmt.Sprintf("Message height: %f\n", height))

	return nil
}
