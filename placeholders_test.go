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
		t, // +2
		strconv.Itoa(line+4),                                             // +3
		callPlaceholderAtDeepStackLevel(PlaceholderLine, LevelDebug, ""), // +4
	) // +5
	assert.Equal( // +6
		t, // +7
		strconv.Itoa(line+9),                                               // +8
		callPlaceholderAtDeepStackLevel(PlaceholderLine, LevelWarning, ""), // +9
	)
}

func TestPlaceholderFile_ReturnsCallerFilenameInShortModeByDefault(
	t *testing.T,
) {
	_, file, _, _ := runtime.Caller(0)

	assert.Equal(
		t,
		filepath.Base(file),
		callPlaceholderAtDeepStackLevel(PlaceholderFile, LevelInfo, ""),
	)
}

func TestPlaceholderFile_ReturnsCallerFilenameInShortModeIfValueIsShort(
	t *testing.T,
) {
	_, file, _, _ := runtime.Caller(0)

	assert.Equal(
		t,
		filepath.Base(file),
		callPlaceholderAtDeepStackLevel(PlaceholderFile, LevelDebug, "short"),
	)
}

func TestPlaceholderFile_ReturnsCallerFullFilenameInLongMode(
	t *testing.T,
) {
	_, file, _, _ := runtime.Caller(0)

	assert.Equal(
		t,
		file,
		callPlaceholderAtDeepStackLevel(PlaceholderFile, LevelDebug, "long"),
	)
	assert.Equal(
		t,
		file,
		callPlaceholderAtDeepStackLevel(PlaceholderFile, LevelWarning, "long"),
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

// Placeholders which uses runtime.Caller for receiving data about caller
// should be tested using this helper because Placeholder function will be
// executed by Formatter at N level of stack trace below than calling some log
// function.
//
// For example:
// 1 level: github.com/user/project: call Log.Warning
// 2 level: lorg: Log.Warning call log.log
// 3 level: lorg: Log.log call Log.doLog
// 4 level: lorg: log.doLog call Format.Render
// 5 level: lorg: Format.Render call Placeholder
//
// Actually, this function is 2nd level of stack trace to call placeholder
// 1st level of stack trace is testcase
func callPlaceholderAtDeepStackLevel(
	placeholder Placeholder, logLevel Level, placeholderValue string,
) string {
	actualStackLevel := 2

	return doRecursiveCallPlaceholder(
		placeholder, logLevel, placeholderValue,
		actualStackLevel, PlaceholderCallStackLevel,
	)
}

func doRecursiveCallPlaceholder(
	placeholder Placeholder, logLevel Level, placeholderValue string,
	actualStackLevel, expectedStackLevel int,
) string {
	actualStackLevel++

	if actualStackLevel >= expectedStackLevel {
		return placeholder(logLevel, placeholderValue)
	}

	return doRecursiveCallPlaceholder(
		placeholder, logLevel, placeholderValue,
		actualStackLevel, expectedStackLevel,
	)
}
