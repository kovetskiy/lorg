package lorg

import (
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

func placeholderLine(logLevel Level, arg string) string {
	_, _, line, ok := runtime.Caller(1)
	if !ok {
		return "??"
	}

	return strconv.Itoa(line)
}
