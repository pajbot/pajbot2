package modules

import (
	"github.com/pajbot/pajbot2/pkg"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
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
	mbase.Base
}

func newTest(b mbase.Base) pkg.Module {
	return &test{
		Base: b,
	}
}
