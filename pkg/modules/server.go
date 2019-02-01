package modules

import (
	"database/sql"
	"strings"
	"sync"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/common/config"
	"github.com/pajlada/pajbot2/pkg/report"
)

type server struct {
	sql          *sql.DB
	oldSession   *sql.DB
	pubSub       pkg.PubSub
	reportHolder *report.Holder
}

var _server server

func InitServer(app pkg.Application, pajbot1Config *config.Pajbot1Config, reportHolder *report.Holder) error {
	var err error

	_server.sql = app.SQL()
	_server.oldSession, err = sql.Open("mysql", pajbot1Config.SQL.DSN)
	_server.pubSub = app.PubSub()
	_server.reportHolder = reportHolder
	if err != nil {
		return err
	}

	return nil
}

type moduleParameterSpec struct {
	description   string
	parameterType string
	defaultValue  interface{}
}

type moduleSpec struct {
	maker pkg.ModuleMaker

	// i.e. "report". This is used in external calls enabling or disabling the module
	// the ID is also what's used when storing settings in the database
	id string

	// i.e. "Report"
	name string

	enabledByDefault bool

	priority int

	parameters map[string]*moduleParameterSpec
}

func (s *moduleSpec) ID() string {
	return s.id
}

func (s *moduleSpec) Name() string {
	return s.name
}

func (s *moduleSpec) EnabledByDefault() bool {
	return s.enabledByDefault
}

func (s *moduleSpec) Maker() pkg.ModuleMaker {
	return s.maker
}

func (s *moduleSpec) Priority() int {
	return s.priority
}

var _ pkg.ModuleSpec = &moduleSpec{}

var _modulesMutex sync.Mutex
var _modules []*moduleSpec

var _validModulesMutex sync.Mutex
var _validModules map[string]*moduleSpec

func Register(spec *moduleSpec) {
	if spec == nil {
		panic("Trying to register a nil module spec")
	}

	if spec.ID() == "" {
		panic("Missing ID in module spec")
	}

	if spec.Name() == "" {
		panic("Missing Name in module spec")
	}

	if spec.Maker() == nil {
		panic("Missing Maker in module spec")
	}

	_modulesMutex.Lock()
	_modules = append(_modules, spec)
	_modulesMutex.Unlock()

	_validModulesMutex.Lock()
	if _validModules == nil {
		_validModules = make(map[string]*moduleSpec)
	}
	_validModules[strings.ToLower(spec.ID())] = spec
	_validModulesMutex.Unlock()
}

func Modules() []*moduleSpec {
	_modulesMutex.Lock()
	defer _modulesMutex.Unlock()

	return _modules
}

func GetModule(moduleID string) (pkg.ModuleSpec, bool) {
	_validModulesMutex.Lock()
	defer _validModulesMutex.Unlock()

	spec, ok := _validModules[strings.ToLower(moduleID)]
	return spec, ok
}

func init() {
	// TODO: Remove action performer
	Register(bttvEmoteParserSpec)

	Register(&badCharacterSpec)
	Register(&bannedNamesSpec)
	Register(&pajbot1BanphraseSpec)
	// TODO: Remove bttv emote parser. This should be done automatically, always
	// custom commands
	Register(&emoteLimitSpec)
	Register(&giveawaySpec)
	Register(&latinFilterSpec)
	Register(&linkFilterSpec)
	Register(&messageLengthLimitSpec)
	Register(&nukeSpec)
	Register(&pajbot1CommandsSpec)
	Register(&reportSpec)
	Register(&testSpec)
	Register(basicCommandsModuleSpec)
	Register(otherCommandsModuleSpec)
	Register(actionPerformerModuleSpec)
}
