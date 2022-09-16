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
			LevelFatal:   {stderr},
			LevelError:   {stderr},
			LevelWarning: {stderr},
			LevelInfo:    {stderr},
			LevelDebug:   {stderr},
			LevelTrace:   {stderr},
		},
		mutex: &sync.Mutex{},
	}
}

func (output *Output) SetLevelWriterCondition(
	level Level, writer ...io.Writer,
) *Output {
	output.mutex.Lock()

	output.conditions[level] = writer

	output.mutex.Unlock()

	return output
}

func (output *Output) Write(buffer []byte) (int, error) {
	panic("should be not called")
}

func (output *Output) WriteWithLevel(
	data []byte, level Level,
) (int, error) {
	output.mutex.Lock()

	writers, ok := output.conditions[level]
	if !ok {
		output.mutex.Unlock()
		return 0, fmt.Errorf("there is no writers for level %s", level)
	}

	var written int
	var err error
	for _, writer := range writers {
		written, err = writer.Write(data)
	}

	output.mutex.Unlock()

	return written, err
}
