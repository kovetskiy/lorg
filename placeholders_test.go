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

func TestPlaceholderLevel_ReturnsLevelStringRepresentation(t *testing.T) {
	assert.Equal(t, "DEBUG", PlaceholderLevel(LevelDebug, "blah"))
	assert.Equal(t, "FATAL", PlaceholderLevel(LevelFatal, "blah"))
	assert.Equal(t, "FATAL", PlaceholderLevel(LevelFatal, ""))
}

func TestPlaceholderLine_ReturnsCallerLine(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	assert.Equal( // +1
		t, strconv.Itoa(line+2), PlaceholderLine(LevelDebug, ""), // +2
	) // +3
	assert.Equal( // +4
		t, strconv.Itoa(line+5), PlaceholderLine(LevelWarning, "blah"), // +5
	)
}

func TestPlaceholderFile_ReturnsCallerFilenameInShortModeByDefault(
	t *testing.T,
) {
	_, file, _, _ := runtime.Caller(0)

	assert.Equal(
		t, filepath.Base(file), PlaceholderFile(LevelInfo, ""),
	)
}

func TestPlaceholderFile_ReturnsCallerFilenameInShortModeIfValueIsShort(
	t *testing.T,
) {
	_, file, _, _ := runtime.Caller(0)

	assert.Equal(
		t, filepath.Base(file), PlaceholderFile(LevelDebug, "short"),
	)
}

func TestPlaceholderFile_ReturnsCallerFullFilenameInLongMode(
	t *testing.T,
) {
	_, file, _, _ := runtime.Caller(0)

	assert.Equal(
		t, file, PlaceholderFile(LevelDebug, "long"),
	)
	assert.Equal(
		t, file, PlaceholderFile(LevelWarning, "long"),
	)
}

func TestPlaceholderTime_ReturnsTimestampIfValueIsTimestamp(t *testing.T) {
	assert.Equal(
		t,
		fmt.Sprint(time.Now().Unix()),
		PlaceholderTime(LevelDebug, "timestamp"),
	)
}

func TestPlaceholderTime_ReturnsTimeUsingDefaultLayoutIfValueNotSpecified(
	t *testing.T,
) {
	assert.Equal(
		t,
		time.Now().Format(PlaceholderTimeDefaultLayout),
		PlaceholderTime(LevelDebug, ""),
	)
}

func TestPlaceholderTime_ReturnsTimeUsingSpecifiedLayout(t *testing.T) {
	assert.Equal(
		t,
		time.Now().Format(time.Kitchen),
		PlaceholderTime(LevelDebug, time.Kitchen),
	)
}
