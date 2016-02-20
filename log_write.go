package lorg

import (
	"fmt"
	"os"
	"strings"
)

func (log *Log) write(message string) error {
	_, err := log.output.Write([]byte(message))
	return err
}

func (log *Log) log(level Level, args []interface{}) {
	if log.level < level {
		return
	}

	format := log.format.Render(level)

	// there is no need for Sprintf, so just replace %s to message
	message := strings.Replace(format, "%s", fmt.Sprint(args...), 1)
	message = message + "\n"

	err := log.write(message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write to log: %#v", err)
	}
}

func (log *Log) logf(
	level Level, format string, args []interface{},
) {
	if log.level < level {
		return
	}

	log.log(
		level,
		[]interface{}{
			fmt.Sprintf(format, args...),
		},
	)
}
