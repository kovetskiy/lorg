package lorg

import (
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
)

var (
	rePlaceholder = regexp.MustCompile(`\${(\w+)(:([^}]+))?}`)

	defaultPlaceholders = map[string]Placeholder{
		"level": placeholderLevel,
		"line":  placeholderLine,
	}
)

type Placeholder func(level Level, arg string) string

func placeholderLevel(logLevel Level, arg string) string {
	return logLevel.String()
}

func placeholderLine(logLevel Level, _ string) string {
	_, _, line, ok := runtime.Caller(1)
	if !ok {
		return "??"
	}

	return strconv.Itoa(line)
}

func placeholderFile(logLevel Level, mode string) string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "??"
	}

	if mode == "long" {
		return file
	}

	return filepath.Base(file)
}
