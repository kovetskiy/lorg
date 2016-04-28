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

func TestPlaceholderLevel_ReturnsLevelWithoutOptions(t *testing.T) {
	test := assert.New(t)

	test.Equal("DEBUG", PlaceholderLevel(LevelDebug, ""))
	test.Equal("FATAL", PlaceholderLevel(LevelFatal, ""))
}

func TestPlaceholderLevel_ReturnsFormattedLevel(t *testing.T) {
	test := assert.New(t)

	test.Equal("DEBUG", PlaceholderLevel(LevelDebug, "%s"))
	test.Equal("[FATAL]", PlaceholderLevel(LevelFatal, "[%s]"))
	test.Equal("xxx INFO xxx", PlaceholderLevel(LevelInfo, "xxx %s xxx"))
}

func TestPlaceholderLevel_ReturnsLevelAlignedLeft(t *testing.T) {
	test := assert.New(t)

	test.Equal("DEBUG  ", PlaceholderLevel(LevelDebug, "%s:left"))
	test.Equal("WARNING", PlaceholderLevel(LevelWarning, "%s:left"))
}

func TestPlaceholderLevel_ReturnsLevelAlignedRight(t *testing.T) {
	test := assert.New(t)

	test.Equal("  DEBUG", PlaceholderLevel(LevelDebug, "%s:right"))
	test.Equal("WARNING", PlaceholderLevel(LevelWarning, "%s:right"))
}

func TestPlaceholderLevel_ReturnsFormattedLevelAlignedLeft(t *testing.T) {
	test := assert.New(t)

	test.Equal("[DEBUG]  ", PlaceholderLevel(LevelDebug, "[%s]:left"))
	test.Equal("[WARNING]", PlaceholderLevel(LevelWarning, "[%s]:left"))
}

func TestPlaceholderLevel_ReturnsFormattedLevelAlignedRight(t *testing.T) {
	test := assert.New(t)

	test.Equal("  [DEBUG]", PlaceholderLevel(LevelDebug, "[%s]:right"))
	test.Equal("[WARNING]", PlaceholderLevel(LevelWarning, "[%s]:right"))
}

func TestPlaceholderLevel_ReturnsLevelAlignedLeftShortString(t *testing.T) {
	test := assert.New(t)

	test.Equal("DEBUG", PlaceholderLevel(LevelDebug, "%s:left:true"))
	test.Equal("WARN ", PlaceholderLevel(LevelWarning, "%s:left:true"))
}

func TestPlaceholderLevel_ReturnsLevelAlignedRightShortString(t *testing.T) {
	test := assert.New(t)

	test.Equal("DEBUG", PlaceholderLevel(LevelDebug, "%s:right:true"))
	test.Equal(" WARN", PlaceholderLevel(LevelWarning, "%s:right:true"))
}

func TestPlaceholderLevel_ReturnsFormattedLevelAlignedLeftShortString(
	t *testing.T,
) {
	test := assert.New(t)

	test.Equal("[DEBUG]", PlaceholderLevel(LevelDebug, "[%s]:left:true"))
	test.Equal("[WARN] ", PlaceholderLevel(LevelWarning, "[%s]:left:true"))
}

func TestPlaceholderLevel_ReturnsFormattedLevelAlignedRightShortString(
	t *testing.T,
) {
	test := assert.New(t)

	test.Equal("[DEBUG]", PlaceholderLevel(LevelDebug, "[%s]:right:true"))
	test.Equal(" [WARN]", PlaceholderLevel(LevelWarning, "[%s]:right:true"))
}

func TestPlaceholderLine_ReturnsCallerLine(t *testing.T) {
	test := assert.New(t)

	_, _, line, _ := runtime.Caller(0)
	test.Equal( // +1
		strconv.Itoa(line+3),                                  // +2
		fakePlaceholderStack(PlaceholderLine, LevelDebug, ""), // +3
	) // +4
	test.Equal( // +5
		strconv.Itoa(line+7),                                    // +6
		fakePlaceholderStack(PlaceholderLine, LevelWarning, ""), // +7
	)
}

func TestPlaceholderFile_ReturnsCallerFilenameInShortModeByDefault(
	t *testing.T,
) {
	test := assert.New(t)

	_, file, _, _ := runtime.Caller(0)

	test.Equal(
		filepath.Base(file),
		fakePlaceholderStack(PlaceholderFile, LevelInfo, ""),
	)
}

func TestPlaceholderFile_ReturnsCallerFilenameInShortModeIfValueIsShort(
	t *testing.T,
) {
	test := assert.New(t)

	_, file, _, _ := runtime.Caller(0)

	test.Equal(
		filepath.Base(file),
		fakePlaceholderStack(PlaceholderFile, LevelDebug, "short"),
	)
}

func TestPlaceholderFile_ReturnsCallerFullFilenameInLongMode(
	t *testing.T,
) {
	test := assert.New(t)

	_, file, _, _ := runtime.Caller(0)

	test.Equal(
		file,
		fakePlaceholderStack(PlaceholderFile, LevelDebug, "long"),
	)
	test.Equal(
		file,
		fakePlaceholderStack(PlaceholderFile, LevelWarning, "long"),
	)
}

func TestPlaceholderTime_ReturnsTimestampIfValueIsTimestamp(t *testing.T) {
	test := assert.New(t)

	test.Equal(
		fmt.Sprint(time.Now().Unix()),
		PlaceholderTime(LevelDebug, "timestamp"),
	)
}

func TestPlaceholderTime_ReturnsTimeUsingDefaultLayoutIfValueNotSpecified(
	t *testing.T,
) {
	test := assert.New(t)

	test.Equal(
		time.Now().Format(PlaceholderTimeDefaultLayout),
		PlaceholderTime(LevelDebug, ""),
	)
}

func TestPlaceholderTime_ReturnsTimeUsingSpecifiedLayout(t *testing.T) {
	test := assert.New(t)

	test.Equal(
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
func fakePlaceholderStack(
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
