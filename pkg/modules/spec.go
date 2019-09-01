package modules

import (
	"github.com/pajbot/pajbot2/pkg"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
)

type moduleParameterSpec struct {
	description   string
	parameterType string
	defaultValue  interface{}
}

type moduleMaker func(b mbase.Base) pkg.Module

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
	b := mbase.New(s, bot, _server.sql, _server.oldSession, _server.pubSub, _server.reportHolder)
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
