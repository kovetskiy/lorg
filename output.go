package lorg

import (
	"fmt"
	"io"
	"sync"
)

type SmartOutput interface {
	io.Writer
	WriteWithLevel([]byte, Level) (int, error)
}

type Output struct {
	conditions map[Level][]io.Writer
	mutex      *sync.Mutex
}

func NewOutput(stderr io.Writer) *Output {
	return &Output{
		conditions: map[Level][]io.Writer{
			LevelFatal:   []io.Writer{stderr},
			LevelError:   []io.Writer{stderr},
			LevelWarning: []io.Writer{stderr},
			LevelInfo:    []io.Writer{stderr},
			LevelDebug:   []io.Writer{stderr},
			LevelTrace:   []io.Writer{stderr},
		},
		mutex: &sync.Mutex{},
	}
}

func (output *Output) SetLevelWriterCondition(
	level Level, writer ...io.Writer,
) *Output {
	output.mutex.Lock()
	defer output.mutex.Unlock()

	output.conditions[level] = writer

	return output
}

func (output *Output) Write(buffer []byte) (int, error) {
	panic("should be not called")
}

func (output *Output) WriteWithLevel(
	data []byte, level Level,
) (int, error) {
	output.mutex.Lock()
	defer output.mutex.Unlock()

	writers, ok := output.conditions[level]
	if !ok {
		return 0, fmt.Errorf("there is no writers for level %s", level)
	}

	var written int
	var err error
	for _, writer := range writers {
		written, err = writer.Write(data)
	}

	return written, err
}
