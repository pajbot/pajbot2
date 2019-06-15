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

	// Returns the bot channel that the module has saved
	BotChannel() BotChannel

	OnWhisper(event MessageEvent) Actions
	OnMessage(event MessageEvent) Actions

	// Implemented in base module
	ID() string
	Type() ModuleType
	Priority() int
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
	Parameters() map[string]ModuleParameterSpec

	Create(bot BotChannel) Module

	Priority() int
}

type ModuleParameterSpec func() ModuleParameter

type ModuleParameter interface {
	Description() string
	DefaultValue() interface{}
}
