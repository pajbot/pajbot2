package commandlist

import (
	"sync"

	"github.com/pajlada/pajbot2/pkg"
)

var (
	commandsMutex sync.Mutex
	commands      []pkg.CommandInfo
)

func Register(commandInfo pkg.CommandInfo) {
	commandsMutex.Lock()
	commands = append(commands, commandInfo)
	commandsMutex.Unlock()
}

func Commands() []pkg.CommandInfo {
	commandsMutex.Lock()
	defer commandsMutex.Unlock()
	return commands
}

func init() {

}
