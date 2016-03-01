package lorg

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

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

func TestPlaceholderFile(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)

	assert.Equal(
		t, filepath.Base(file), placeholderFile(LevelDebug, "short"),
	)
	assert.Equal(
		t, filepath.Base(file), placeholderFile(LevelInfo, ""),
	)

	assert.Equal(
		t, file, placeholderFile(LevelDebug, "long"),
	)
	assert.Equal(
		t, file, placeholderFile(LevelWarning, "long"),
	)
}

func TestPlaceholderTime(t *testing.T) {
	assert.Equal(
		t,
		fmt.Sprint(time.Now().Unix()),
		placeholderTime(LevelDebug, "timestamp"),
	)

	assert.Equal(
		t,
		time.Now().Format(placeholderTimeLayout),
		placeholderTime(LevelDebug, ""),
	)

	assert.Equal(
		t,
		time.Now().Format(time.Kitchen),
		placeholderTime(LevelDebug, time.Kitchen),
	)
}
