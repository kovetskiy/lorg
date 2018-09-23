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

	text := fmt.Sprint(value...)

	shift := log.shiftIndent
	if shift == 0 && log.indentLines {
		shift = strings.Index(format, "%s")
	}

	if shift > 0 {
		text = indent(text, shift)
	}

	// here is no need for Sprintf, so just replace %s to text
	entry := strings.Replace(format, "%s", text, 1) + "\n"

	func() {
		log.mutex.Lock()
		defer log.mutex.Unlock()

		err := log.write(entry, level)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write to log: %#v", err)
		}
	}()
}

func (log *Log) write(text string, level Level) error {
	_, err := log.output.WriteWithLevel([]byte(text), level)
	return err
}

func indent(text string, shift int) string {
	text = strings.Replace(
		text,
		"\n",
		"\n"+strings.Repeat(" ", shift),
		-1,
	)

	return text
}
