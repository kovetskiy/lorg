package lorg

import (
	"io"
	"os"
)

const (
	DefaultFormatting = `${date} ${level} %s`
)

var (
	defaultLevel  = LevelInfo
	defaultFormat = NewFormat(DefaultFormatting)
	defaultOutput = os.Stderr
)

type Log struct {
	level  Level
	output io.Writer
	format Formatter
}

func NewLog() *Log {
	log := &Log{
		level:  defaultLevel,
		format: defaultFormat,
		output: defaultOutput,
	}

	return log
}

func (log *Log) SetLevel(level Level) {
	log.level = level
}

func (log *Log) SetFormat(format Formatter) {
	log.format = format
}

func (log *Log) SetOutput(output io.Writer) {
	log.output = output
}

func (log *Log) Fatal(args ...interface{}) {
	log.log(LevelFatal, args)
	os.Exit(1)
}

func (log *Log) Fatalf(format string, args ...interface{}) {
	log.logf(LevelFatal, format, args)
	os.Exit(1)
}

func (log *Log) Error(args ...interface{}) {
	log.log(LevelError, args)
}

func (log *Log) Errorf(format string, args ...interface{}) {
	log.logf(LevelError, format, args)
}

func (log *Log) Warning(args ...interface{}) {
	log.log(LevelWarning, args)
}

func (log *Log) Warningf(format string, args ...interface{}) {
	log.logf(LevelWarning, format, args)
}

// Pseudonim for Log.Info
func (log *Log) Print(args ...interface{}) {
	log.log(LevelInfo, args)
}

// Pseudonim for Log.Infof
func (log *Log) Printf(format string, args ...interface{}) {
	log.logf(LevelInfo, format, args)
}

func (log *Log) Info(args ...interface{}) {
	log.log(LevelInfo, args)
}

func (log *Log) Infof(format string, args ...interface{}) {
	log.logf(LevelInfo, format, args)
}

func (log *Log) Debug(args ...interface{}) {
	log.log(LevelDebug, args)
}

func (log *Log) Debugf(format string, args ...interface{}) {
	log.logf(LevelDebug, format, args)
}
