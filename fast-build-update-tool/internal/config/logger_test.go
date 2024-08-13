package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeLoggerVerbose(t *testing.T) {
	l, err := InitializeLogger(true)
	assert.Nil(t, err)
	l.Close()
}

func TestInitializeLoggerNotVerbose(t *testing.T) {
	l, err := InitializeLogger(false)
	assert.Nil(t, err)
	l.Close()
}
