package modules

import (
	"github.com/pajbot/pajbot2/pkg"
)

func init() {
	Register("test", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:               "test",
			name:             "Test",
			enabledByDefault: false,

			maker: newTest,
		}
	})
}

type test struct {
	base
}

func newTest(b base) pkg.Module {
	return &test{
		base: b,
	}
}
