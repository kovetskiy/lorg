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
	assert.Equal(t, "DEBUG", PlaceholderLevel(LevelDebug, "blah"))
	assert.Equal(t, "FATAL", PlaceholderLevel(LevelFatal, "blah"))
	assert.Equal(t, "FATAL", PlaceholderLevel(LevelFatal, ""))
}

func TestPlaceholderLine(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	assert.Equal( // +1
		t, strconv.Itoa(line+2), PlaceholderLine(LevelDebug, ""), // +2
	) // +3
	assert.Equal( // +4
		t, strconv.Itoa(line+5), PlaceholderLine(LevelWarning, "blah"), // +5
	)
}

func TestPlaceholderFile(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)

	assert.Equal(
		t, filepath.Base(file), PlaceholderFile(LevelDebug, "short"),
	)
	assert.Equal(
		t, filepath.Base(file), PlaceholderFile(LevelInfo, ""),
	)

	assert.Equal(
		t, file, PlaceholderFile(LevelDebug, "long"),
	)
	assert.Equal(
		t, file, PlaceholderFile(LevelWarning, "long"),
	)
}

func TestPlaceholderTime(t *testing.T) {
	assert.Equal(
		t,
		fmt.Sprint(time.Now().Unix()),
		PlaceholderTime(LevelDebug, "timestamp"),
	)

	assert.Equal(
		t,
		time.Now().Format(placeholderTimeLayout),
		PlaceholderTime(LevelDebug, ""),
	)

	assert.Equal(
		t,
		time.Now().Format(time.Kitchen),
		PlaceholderTime(LevelDebug, time.Kitchen),
	)
}
