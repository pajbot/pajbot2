package modules

import (
	"database/sql"
	"sync"

	"github.com/garyburd/redigo/redis"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/common/config"
	"github.com/pajlada/pajbot2/pkg/pubsub"
	"github.com/pajlada/pajbot2/pkg/report"
)

type server struct {
	redis        *redis.Pool
	sql          *sql.DB
	oldSession   *sql.DB
	pubSub       *pubsub.PubSub
	reportHolder *report.Holder
}

var _server server

func InitServer(redis *redis.Pool, _sql *sql.DB, pajbot1Config config.Pajbot1Config, pubSub *pubsub.PubSub, reportHolder *report.Holder) error {
	var err error

	_server.redis = redis
	_server.sql = _sql
	_server.oldSession, err = sql.Open("mysql", pajbot1Config.SQL.DSN)
	_server.pubSub = pubSub
	_server.reportHolder = reportHolder
	if err != nil {
		return err
	}

	return nil
}

type moduleSpec struct {
	maker pkg.ModuleMaker

	// i.e. "report". This is used in external calls enabling or disabling the module
	// the ID is also what's used when storing settings in the database
	id string

	// i.e. "Report"
	name string

	enabledByDefault bool

	Priority int
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

var _ pkg.ModuleSpec = &moduleSpec{}

var _modulesMutex sync.Mutex
var _modules []*moduleSpec

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
	defer _modulesMutex.Unlock()
	_modules = append(_modules, spec)
}

func Modules() []*moduleSpec {
	_modulesMutex.Lock()
	defer _modulesMutex.Unlock()

	return _modules
}

func init() {
	// TODO: Remove action performer
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
	// Register(&messageHeightLimitSpec)
	Register(&nukeSpec)
	Register(&pajbot1CommandsSpec)
	Register(&reportSpec)
	Register(&testSpec)
	Register(basicCommandsModuleSpec)
}
