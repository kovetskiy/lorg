package lorg

import "regexp"

var (
	rePlaceholder = regexp.MustCompile(`\${(\w+)(:([^}]+))?}`)

	defaultPlaceholders = map[string]Placeholder{
		"level": placeholderLevel,
	}
)

type Placeholder func(level Level, arg string) string

func placeholderLevel(logLevel Level, arg string) string {
	return logLevel.String()
}
