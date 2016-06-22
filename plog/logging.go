package plog

import (
	"os"

	"github.com/op/go-logging"
)

var (
	log    = logging.MustGetLogger("pajbot2")
	format = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000} %{shortpkg:-12s} %{shortfile:-19s} %{level:.4s}%{color:reset} %{message}`,
	)
)

// InitLogging xD
func InitLogging() {
	backend1 := logging.NewLogBackend(os.Stdout, "", 0)
	backend2 := logging.NewLogBackend(os.Stdout, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled, backend2Formatter)
}

// GetLogger xD
func GetLogger() *logging.Logger {
	return log
}
