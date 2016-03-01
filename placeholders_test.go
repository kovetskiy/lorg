package lorg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlaceholderLevel(t *testing.T) {
	assert.Equal(t, "DEBUG", placeholderLevel(LevelDebug, "blah"))
	assert.Equal(t, "FATAL", placeholderLevel(LevelFatal, "blah"))
	assert.Equal(t, "FATAL", placeholderLevel(LevelFatal, ""))
}
