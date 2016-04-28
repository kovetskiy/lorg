package lorg

import (
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/zazab/zhash"
)

// Placeholder is function which will be called by Formatter for all parsed
// placeholders.
//
// Actually, placeholder is just a string which will be replaced onto result
// of specified Placeholder function.
//
// logLevel is level of current logging record, value is optional parameter
// which can be passed from Formatter if placeholder defined with argument.
// Example:
//     * ${level}       - value will be ""
//     * ${level:test}  - value will be "test"
//     * ${level:a:b:c} - value will be "a:b:c"
type Placeholder func(logLevel Level, value string) string

var (
	rePlaceholder = regexp.MustCompile(`\${(\w+)(:([^}]+))?}`)

	defaultPlaceholders = map[string]Placeholder{
		"level": PlaceholderLevel,
		"line":  PlaceholderLine,
		"file":  PlaceholderFile,
		"time":  PlaceholderTime,
	}

	cache = zhash.NewHash()
)

const (
	// PlaceholderTimeDefaultLayout is the string represents of time layout
	// which will be used by PlaceholderTime as default layout.
	//
	// See more about time layouts in package time documentation.
	PlaceholderTimeDefaultLayout = "2006-01-02 15:04:05"

	// PlaceholderCallStackLevel should be used as argument to
	// runtime.Caller if placeholder want to get information about calling log
	// functions.
	PlaceholderCallStackLevel = 5
)

// PlaceholderLevel returns level of current logging record.
//
// Using: ${level}
func PlaceholderLevel(logLevel Level, optional string) string {
	path := []string{
		"placeholders", "level", logLevel.String(), optional,
	}

	cached, err := cache.GetString(path...)
	if err == nil {
		return cached
	}

	const (
		maxLevelStringLength      = 7
		maxLevelStringLengthShort = 5
	)

	var (
		format    = "%s"
		align     = false
		alignment = "left"
		short     = false
	)

	options := options(optional, 3)

	if options[0] != "" {
		format = options[0]
	}

	if options[1] == "left" || options[1] == "right" {
		alignment = options[1]
		align = true
	}

	if isTrueString(options[2]) {
		short = true
	}

	var levelString string
	if short {
		levelString = logLevel.StringShort()
	} else {
		levelString = logLevel.String()
	}

	value := fmt.Sprintf(format, levelString)

	if align {
		var shift int
		if short {
			shift = maxLevelStringLengthShort - len(levelString)
		} else {
			shift = maxLevelStringLength - len(levelString)
		}

		switch alignment {
		case "left":
			value = value + strings.Repeat(" ", shift)
		case "right":
			value = strings.Repeat(" ", shift) + value
		}
	}

	cache.Set(value, path...)

	return value
}

// PlaceholderLine returns a file line where has been called logging function.
//
// Using: ${line}
func PlaceholderLine(logLevel Level, _ string) string {
	_, _, line, ok := runtime.Caller(PlaceholderCallStackLevel)
	if !ok {
		return "??"
	}

	return strconv.Itoa(line)
}

// PlaceholderFile returns a file name where has been called logging function.
// PlaceholderFile can work in two modes:
//    * "short":   default behaviour, a final file name will be returned after
//                     passing to filepath.Base function.
//                     Using: ${file:short} or just ${file}
//    * "long":    a final file name will be retuned as is.
//                     Using: ${file:long}
func PlaceholderFile(logLevel Level, mode string) string {
	_, file, _, ok := runtime.Caller(PlaceholderCallStackLevel)
	if !ok {
		return "??"
	}

	if mode == "long" {
		return file
	}

	return filepath.Base(file)
}

// PlaceholderTime returns current time formatted with specified time
// formatting layout. If formatting layout is not specified, PlaceholderTime
// will use const PlaceholderTimeDefaultLayout as layout.
//
// PlaceholderTime returns unix timestamp if layout specified as "timestamp"
// Otherwise layout will be passed to time.Time.Format function as is.
//
// Using: ${time}
//        ${time:timestamp}
//        ${time:15:04:05}
func PlaceholderTime(logLevel Level, layout string) string {
	if layout == "timestamp" {
		return fmt.Sprint(time.Now().Unix())
	}

	if layout == "" {
		layout = PlaceholderTimeDefaultLayout
	}

	return time.Now().Format(layout)
}

func isTrueString(str string) bool {
	return str == "true" || str == "yes" || str == "1"
}

func options(str string, count int) []string {
	cached, err := cache.GetStringSlice("options", str, strconv.Itoa(count))
	if err == nil {
		return cached
	}

	options := strings.SplitN(
		strings.Replace(str, `\:`, "\x00", -1), ":", count,
	)

	for index, option := range options {
		options[index] = strings.Replace(option, "\x00", ":", -1)
	}

	for len(options) < count {
		options = append(options, "")
	}

	cache.Set(options, "options", str, strconv.Itoa(count))

	return options
}
