package plog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogging(t *testing.T) {
	// Not sure what we can try here.
	// Just making sure it doesn't crash? KKona
	InitLogging()
}

func TestGetLogger(t *testing.T) {
	assert.NotNil(t, GetLogger())
}
