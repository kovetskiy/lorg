package lorg

import (
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
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
	var (
		format = "%s"
		align  = false
	)

	options := strings.SplitN(optional, ":", 2)
	if options[0] != "" {
		format = options[0]
	}

	if len(options) >= 2 && (options[1] == "true" || options[1] == "1") {
		align = true
	}

	const maxLevelStringLength = 7

	value := fmt.Sprintf(format, logLevel.String())
	if align {
		alignment := maxLevelStringLength - len(logLevel.String())
		if alignment > 0 {
			value += strings.Repeat(" ", alignment)
		}
	}

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
