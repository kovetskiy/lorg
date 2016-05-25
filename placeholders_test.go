package lorg

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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
		strconv.Itoa(line+3),                        // +2
		callLoggerWithFormat("${line}", LevelDebug), // +3
	) // +4
	test.Equal( // +5
		strconv.Itoa(line+7),                          // +6
		callLoggerWithFormat("${line}", LevelWarning), // +7
	)
}

func TestPlaceholderFile_ReturnsCallerFilenameInShortModeByDefault(
	t *testing.T,
) {
	test := assert.New(t)

	_, file, _, _ := runtime.Caller(0)

	test.Equal(
		filepath.Base(file),
		callLoggerWithFormat("${file}", LevelInfo),
	)
}

func TestPlaceholderFile_ReturnsCallerFilenameInShortModeIfValueIsShort(
	t *testing.T,
) {
	test := assert.New(t)

	_, file, _, _ := runtime.Caller(0)

	test.Equal(
		filepath.Base(file),
		callLoggerWithFormat("${file:short}", LevelDebug),
	)
}

func TestPlaceholderFile_ReturnsCallerFullFilenameInLongMode(
	t *testing.T,
) {
	test := assert.New(t)

	_, file, _, _ := runtime.Caller(0)

	test.Equal(
		file,
		callLoggerWithFormat("${file:long}", LevelDebug),
	)
	test.Equal(
		file,
		callLoggerWithFormat("${file:long}", LevelWarning),
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
func callLoggerWithFormat(
	format string, logLevel Level,
) string {
	placeholderCallStackLevel += 1
	defer func() {
		placeholderCallStackLevel = PlaceholderCallStackLevel
	}()

	buffer := bytes.NewBuffer(nil)
	logger := NewLog()
	logger.SetOutput(buffer)
	logger.SetFormat(NewFormat(format))
	logger.SetLevel(LevelDebug)

	switch logLevel {
	case LevelDebug:
		logger.Debug("")
	case LevelWarning:
		logger.Warning("")

	case LevelInfo:
		logger.Info("")

	default:
		panic("unexpected logging level: " + logLevel.String())
	}

	return strings.TrimRight(string(buffer.Bytes()), "\n")
}
