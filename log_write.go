package lorg

import (
	"fmt"
	"os"
	"strings"
)

func (log *Log) log(level Level, value ...interface{}) {
	if log.level < level {
		return
	}

	log.doLog(level, value...)
}

func (log *Log) logf(level Level, format string, value ...interface{}) {
	if log.level < level {
		return
	}

	log.doLog(level, fmt.Sprintf(format, value...))
}

func (log *Log) doLog(level Level, value ...interface{}) {
	format := log.format.Render(level)

	// there is no need for Sprintf, so just replace %s to message
	message := strings.Replace(format, "%s", fmt.Sprint(value...), 1)
	message = message + "\n"

	err := log.write(message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write to log: %#v", err)
	}
}

func (log *Log) write(message string) error {
	_, err := log.output.Write([]byte(message))
	return err
}
