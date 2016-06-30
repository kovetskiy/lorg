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
	var format string

	func() {
		log.mutex.Lock()
		defer log.mutex.Unlock()

		format = log.format.Render(level, log.prefix)
	}()

	// here is no need for Sprintf, so just replace %s to message
	message := strings.Replace(format, "%s", fmt.Sprint(value...), 1)
	message = message + "\n"

	func() {
		log.mutex.Lock()
		defer log.mutex.Unlock()

		err := log.write(message, level)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write to log: %#v", err)
		}
	}()
}

func (log *Log) write(message string, level Level) error {
	_, err := log.output.WriteWithLevel([]byte(message), level)
	return err
}
