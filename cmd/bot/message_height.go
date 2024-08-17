package main

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

	"github.com/pajbot/utils"
)

var clrInitialized = false

var messageHeightLimitLibraryInitialized = false
var charMapPath string

func initCLR() error {
	executableDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	clrPath := utils.GetEnv("LIBCOREFOLDER", "/usr/share/dotnet/shared/Microsoft.NETCore.App/2.1.5")

	// Path to our own executable
	clr1 := C.CString(executableDir + "/bot")

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

func getMessageHeight(channelName string, channelID string, message string, username string, badgeCountIn int) (float32, error) {
	err := initCLR()
	if err != nil {
		return 0, err
	}
	err = initMessageHeightLimitLibrary()
	if err != nil {
		return 0, err
	}

	initChannel(channelName, channelID)

	channelString := C.CString(channelName)
	input := C.CString(message)
	loginName := C.CString(username)
	displayName := C.CString(username) // TODO

	var emoteStrings []*C.char

	var te []C.TwitchEmote

	var pArray unsafe.Pointer

	if len(te) > 0 {
		pArray = unsafe.Pointer(&te[0])
	}

	badgeCount := C.int(badgeCountIn)

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

	return float32(height), nil
}

func initChannel(channelName, channelID string) error {
	channel := C.CString(channelName)
	cChannelID := C.CString(channelID)

	defer func() {
		C.free(unsafe.Pointer(channel))
		C.free(unsafe.Pointer(cChannelID))
	}()

	res := C.InitChannel2(channel, cChannelID, 5000, true)

	if res != 1 {
		return errors.New("Failed to init Channel " + channelName)
	}

	return nil
}

func initMessageHeightLimitLibrary() error {
	charMap := C.CString(charMapPath)

	res := C.InitCharMap(charMap)

	C.free(unsafe.Pointer(charMap))

	if res != 1 {
		return errors.New(fmt.Sprintf("Failed to init CharMap: %d", int(res)))
	}

	messageHeightLimitLibraryInitialized = true

	return nil
}
