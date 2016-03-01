package lorg

import (
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlaceholderLevel(t *testing.T) {
	assert.Equal(t, "DEBUG", placeholderLevel(LevelDebug, "blah"))
	assert.Equal(t, "FATAL", placeholderLevel(LevelFatal, "blah"))
	assert.Equal(t, "FATAL", placeholderLevel(LevelFatal, ""))
}

func TestPlaceholderLine(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	assert.Equal( // +1
		t, strconv.Itoa(line+2), placeholderLine(LevelDebug, ""), // +2
	) // +3
	assert.Equal( // +4
		t, strconv.Itoa(line+5), placeholderLine(LevelWarning, "blah"), // +5
	)
}
