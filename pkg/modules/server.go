package modules

import (
	"database/sql"
	"sync"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/common/config"
	"github.com/pajbot/pajbot2/pkg/report"
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
	if err != nil {
		return err
	}
	_server.pubSub = app.PubSub()
	_server.reportHolder = reportHolder

	return nil
}

type moduleParameterSpec struct {
	description   string
	parameterType string
	defaultValue  interface{}
}

type moduleMaker func(b base) pkg.Module

type moduleSpec struct {
	// i.e. "report". This is used in external calls enabling or disabling the module
	// the ID is also what's used when storing settings in the database
	id string

	// i.e. "Report"
	name string

	moduleType pkg.ModuleType

	enabledByDefault bool

	priority int

	parameters map[string]pkg.ModuleParameterSpec

	maker moduleMaker
}

func (s *moduleSpec) ID() string {
	return s.id
}

func (s *moduleSpec) Name() string {
	return s.name
}

func (s *moduleSpec) Type() pkg.ModuleType {
	return s.moduleType
}

func (s *moduleSpec) EnabledByDefault() bool {
	return s.enabledByDefault
}

func (s *moduleSpec) Create(bot pkg.BotChannel) pkg.Module {
	b := newBase(s, bot)
	m := s.maker(b)

	return m
}

func (s *moduleSpec) Priority() int {
	return s.priority
}

func (s *moduleSpec) Parameters() map[string]pkg.ModuleParameterSpec {
	return s.parameters
}

var _ pkg.ModuleSpec = &moduleSpec{}

var moduleSpecsMutex sync.Mutex
var moduleSpecs []pkg.ModuleSpec
var moduleSpecsMap map[string]pkg.ModuleSpec

var moduleFactoriesMutex sync.Mutex
var moduleFactories map[string]pkg.ModuleFactory

func Register(moduleID string, factory pkg.ModuleFactory) {
	if factory == nil {
		panic("Trying to register a nil factory")
	}

	moduleFactoriesMutex.Lock()
	if moduleFactories == nil {
		moduleFactories = make(map[string]pkg.ModuleFactory)
	}

	if _, ok := moduleFactories[moduleID]; ok {
		panic("A module factory with the id '" + moduleID + "' was already registered")
	}

	moduleFactories[moduleID] = factory
	moduleFactoriesMutex.Unlock()
}

func generateSpecs() (specs []pkg.ModuleSpec) {
	for _, moduleFactory := range moduleFactories {
		spec := moduleFactory()
		specs = append(specs, spec)
	}

	return
}

func initModuleMap(specs []pkg.ModuleSpec) {
	moduleSpecsMap = make(map[string]pkg.ModuleSpec)

	for _, spec := range specs {
		moduleSpecsMap[spec.ID()] = spec
	}
}

func Modules() []pkg.ModuleSpec {
	moduleSpecsMutex.Lock()
	defer moduleSpecsMutex.Unlock()

	if moduleSpecs == nil {
		moduleSpecs = generateSpecs()
		initModuleMap(moduleSpecs)
	}

	return moduleSpecs
}

func ModulesMap() map[string]pkg.ModuleSpec {
	moduleSpecsMutex.Lock()
	defer moduleSpecsMutex.Unlock()

	if moduleSpecs == nil {
		moduleSpecs = generateSpecs()
		initModuleMap(moduleSpecs)
	}

	return moduleSpecsMap
}

func GetModuleSpec(moduleID string) (spec pkg.ModuleSpec, ok bool) {
	moduleSpecsMutex.Lock()
	defer moduleSpecsMutex.Unlock()

	if moduleSpecs == nil {
		moduleSpecs = generateSpecs()
		initModuleMap(moduleSpecs)
	}

	spec, ok = moduleSpecsMap[moduleID]
	return
}

// GetModuleFactory returns the module factory with the given ID.
// This is useful for tests where the module spec should not be shared
func GetModuleFactory(moduleID string) (factory pkg.ModuleFactory, ok bool) {
	moduleFactoriesMutex.Lock()
	defer moduleFactoriesMutex.Unlock()

	factory, ok = moduleFactories[moduleID]
	return
}
