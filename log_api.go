package lorg

import (
	"io"
	"os"
)

const (
	// DefaultFormatting describes default format for log records
	// If you don't run SetFormat, Log instance will use Format instance with
	// this given formatting.
	//
	// See Format structure documentation for information about `${date}` and
	// `${level}` placeholders.
	DefaultFormatting = `${date} ${level} %s`
)

var (
	defaultLevel  = LevelInfo
	defaultFormat = NewFormat(DefaultFormatting)
	defaultOutput = os.Stderr
)

// Log is the actual log which creates log records based on the functins called
// and passes them to the underlying output.
//
// It's not recommended to instantiate Log without using NewLog() function,
// because Log fields can be changed or added other fields, so it can provide
// bugs in future.
type Log struct {
	level  Level
	output io.Writer
	format Formatter
}

// Creates a new Log instance with default configuration:
// *   default logging level is the LevelInfo, which can be changed
//         using log.SetLevel(Level) method
// *   default logging format is the Format{} instance, which implements
//         Formatter interface and can be changed
//         using log.SetFormat(Formatter) method
// *   default logging output is a stderr (os.Stderr), but also can be changed
//         using log.SetOutput(io.Writer) method
func NewLog() *Log {
	log := &Log{
		level:  defaultLevel,
		format: defaultFormat,
		output: defaultOutput,
	}

	return log
}

// SetLevel sets the logging level for the given log.
// After setting level, logger will log records with same level or above.
//
// Running SetLevel it's not required operation, by default Log instance
// creates with INFO level, so levels above (warn, err, fatal, info) will be
// logged also.
func (log *Log) SetLevel(level Level) {
	log.level = level
}

// SetFormat sets the logging format for the given log.
// All log records will be formatted using specified formatter.
//
// Running SetFormat it's not required operation, by default Log instance
// creates with default Format instance, which have format with level and date
// placeholders.
// See: DefaultFormatting
func (log *Log) SetFormat(format Formatter) {
	log.format = format
}

// The log records are fmt.Printfd'd to this io.Writer.
//
// Running SetOutput it's not required operation, by default Log instance
// logs all records to stderr (os.Stderr)
func (log *Log) SetOutput(output io.Writer) {
	log.output = output
}

// Fatal logs record if given logger level is equal or above LevelFatal, and
// calls os.Exit(1) after logging.
// Arguments are handled in the manner of fmt.Print.
func (log *Log) Fatal(args ...interface{}) {
	log.log(LevelFatal, args)
	os.Exit(1)
}

// Debugf logs record if given logger level is equal or above LevelDebug, and
// calls os.Exit(1) after logging.
// Arguments are handled in the manner of fmt.Print.
func (log *Log) Fatalf(format string, args ...interface{}) {
	log.logf(LevelFatal, format, args)
	os.Exit(1)
}

// Error logs record if given logger level is equal or above LevelError.
// Arguments are handled in the manner of fmt.Print.
func (log *Log) Error(args ...interface{}) {
	log.log(LevelError, args)
}

// Errorf logs record if given logger level is equal or above LevelError.
// Arguments are handled in the manner of fmt.Printf.
func (log *Log) Errorf(format string, args ...interface{}) {
	log.logf(LevelError, format, args)
}

// Warning logs record if given logger level is equal or above LevelWarning.
// Arguments are handled in the manner of fmt.Print.
func (log *Log) Warning(args ...interface{}) {
	log.log(LevelWarning, args)
}

// Warningf logs record if given logger level is equal or above LevelWarning.
// Arguments are handled in the manner of fmt.Printf.
func (log *Log) Warningf(format string, args ...interface{}) {
	log.logf(LevelWarning, format, args)
}

// Printf is pseudonym for Info
func (log *Log) Print(args ...interface{}) {
	log.Info(args...)
}

// Printf is pseudonym for Infof
func (log *Log) Printf(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Info logs record if given logger level is equal or above LevelInfo.
// Arguments are handled in the manner of fmt.Print.
func (log *Log) Info(args ...interface{}) {
	log.log(LevelInfo, args)
}

// Info logs record if given logger level is equal or above LevelInfo.
// Arguments are handled in the manner of fmt.Printf.
func (log *Log) Infof(format string, args ...interface{}) {
	log.logf(LevelInfo, format, args)
}

// Debug logs record if given logger level is equal or above LevelDebug.
// Arguments are handled in the manner of fmt.Print.
func (log *Log) Debug(args ...interface{}) {
	log.log(LevelDebug, args)
}

// Debug logs record if given logger level is equal or above LevelDebug.
// Arguments are handled in the manner of fmt.Printf.
func (log *Log) Debugf(format string, args ...interface{}) {
	log.logf(LevelDebug, format, args)
}
