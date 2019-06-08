package pkg

type ModuleFactory func() ModuleSpec

// A module is local to a bots channel
// i.e. bot "pajbot" joins channels "pajlada" and "forsen"
// Module list looks like this:
// "pajbot":
//  - "pajlada":
//		- "MyTestModule"
//		- "MyTestModule2"
//  - "forsen":
//		- "MyTestModule"
type Module interface {
	// After the module struct is created, it must be initialized with the channel
	LoadSettings([]byte) error

	// Called when the module is disabled. The module can do any cleanup it needs to do here
	Disable() error

	// Returns the spec for the module
	Spec() ModuleSpec

	// Returns the bot channel that the module has saved
	BotChannel() BotChannel

	OnWhisper(bot BotChannel, user User, message Message) error
	OnMessage(bot BotChannel, user User, message Message, action Action) error
}

type ModuleType uint

const (
	ModuleTypeUnsorted = 0
	ModuleTypeFilter   = 1
)

type ModuleSpec interface {
	ID() string
	Name() string
	Type() ModuleType
	EnabledByDefault() bool

	Create(bot BotChannel) Module

	Priority() int
}
