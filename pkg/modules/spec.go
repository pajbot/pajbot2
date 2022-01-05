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

type moduleMaker func(b *mbase.Base) pkg.Module

type Spec struct {
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

type Option func(s *Spec)

func WithModuleType(moduleType pkg.ModuleType) Option {
	return func(s *Spec) {
		s.moduleType = moduleType
	}
}

func NewSpec(id, name string, enabledByDefault bool, maker moduleMaker, opts ...Option) *Spec {
	s := &Spec{
		id:               id,
		name:             name,
		enabledByDefault: enabledByDefault,
		maker:            maker,
	}

	for _, option := range opts {
		option(s)
	}

	return s
}

func (s *Spec) ID() string {
	return s.id
}

func (s *Spec) Name() string {
	return s.name
}

func (s *Spec) Type() pkg.ModuleType {
	return s.moduleType
}

func (s *Spec) EnabledByDefault() bool {
	return s.enabledByDefault
}

func (s *Spec) Create(bot pkg.BotChannel) pkg.Module {
	b := mbase.New(s, bot, _server.sql, _server.oldSession, _server.pubSub, _server.reportHolder)
	m := s.maker(&b)

	return m
}

func (s *Spec) Priority() int {
	return s.priority
}

func (s *Spec) Parameters() map[string]pkg.ModuleParameterSpec {
	return s.parameters
}

func (s *Spec) SetParameters(parameters map[string]pkg.ModuleParameterSpec) {
	s.parameters = parameters
}

var _ pkg.ModuleSpec = &Spec{}
