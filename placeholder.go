package lorg

import (
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

type Placeholder func(level Level, arg string) string

var (
	rePlaceholder = regexp.MustCompile(`\${(\w+)(:([^}]+))?}`)

	defaultPlaceholders = map[string]Placeholder{
		"level": PlaceholderLevel,
		"line":  PlaceholderLine,
		"file":  PlaceholderFile,
		"time":  PlaceholderTime,
	}
)

const (
	placeholderTimeLayout = "2006-01-02 15:04:05"
)

func PlaceholderLevel(logLevel Level, arg string) string {
	return logLevel.String()
}

func PlaceholderLine(logLevel Level, _ string) string {
	_, _, line, ok := runtime.Caller(1)
	if !ok {
		return "??"
	}

	return strconv.Itoa(line)
}

func PlaceholderFile(logLevel Level, mode string) string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "??"
	}

	if mode == "long" {
		return file
	}

	return filepath.Base(file)
}

func PlaceholderTime(logLevel Level, layout string) string {
	if layout == "timestamp" {
		return fmt.Sprint(time.Now().Unix())
	}

	if layout == "" {
		layout = placeholderTimeLayout
	}

	return time.Now().Format(layout)
}
