package lorg

import (
	"io"
	"os"
	"sync"
)

const (
	// DefaultFormatting describes default format for log records
	// If you don't run SetFormat, Log instance will use Format instance with
	// this given formatting.
	//
	// See Format structure documentation for information about `${date}` and
	// `${level}` placeholders.
	DefaultFormatting = `${time} ${level:[%s]\::right:true} ${prefix}%s`
)

var (
	defaultLevel  = LevelInfo
	defaultFormat = NewFormat(DefaultFormatting)
	defaultOutput = NewOutput(os.Stderr)

	// Exiter will be called after Fatal/Fatalf invocation.
	Exiter = os.Exit
)

// Log is the actual log which creates log records based on the functins called
// and passes them to the underlying output.
//
// It's not recommended to instantiate Log without using NewLog() function,
// because Log fields can be changed or added other fields, so it can provide
// bugs in future.
type Log struct {
	level Level

	output      SmartOutput
	format      Formatter
	indentLines bool
	mutex       *sync.Mutex
	children    []*Log
	prefix      string
}

// NewLog creates a new Log instance with default configuration:
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
		mutex:  &sync.Mutex{},
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
	log.mutex.Lock()
	defer log.mutex.Unlock()

	log.level = level

	for _, child := range log.children {
		child.SetLevel(level)
	}
}

// SetFormat sets the logging format for the given log.
// All log records will be formatted using specified formatter.
//
// Running SetFormat it's not required operation, by default Log instance
// creates with default Format instance, which have format with level and date
// placeholders.
// See: DefaultFormatting
func (log *Log) SetFormat(format Formatter) {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	log.format = format
}

// SetOutput sets output of given log instance, the log records are
// fmt.Printfd'd to specified io.Writer.
//
// Running SetOutput it's not required operation, by default Log instance
// logs all records to stderr (os.Stderr)
func (log *Log) SetOutput(output io.Writer) {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	if _, ok := output.(SmartOutput); !ok {
		output = NewOutput(output)
	}

	log.output = output.(SmartOutput)
}

// SetIndentLines changes Log's option that responsible for indenting log entry
// lines in one format.
// With this option log entries with newline symbols will be indented like as
// following:
//
// [INFO] before-new-line
//        after-new-line
//
func (log *Log) SetIndentLines(value bool) {
	log.indentLines = value
}

// Fatal logs record if given logger level is equal or above LevelFatal, and
// calls os.Exit(1) after logging.
// Arguments are handled in the manner of fmt.Print.
func (log *Log) Fatal(value ...interface{}) {
	log.log(LevelFatal, value...)
	Exiter(1)
}

// Fatalf logs record if given logger level is equal or above LevelFatal, and
// calls os.Exit(1) after logging.
// Arguments are handled in the manner of fmt.Print.
func (log *Log) Fatalf(format string, value ...interface{}) {
	log.logf(LevelFatal, format, value...)
	Exiter(1)
}

// Error logs record if given logger level is equal or above LevelError.
// Arguments are handled in the manner of fmt.Print.
func (log *Log) Error(value ...interface{}) {
	log.log(LevelError, value...)
}

// Errorf logs record if given logger level is equal or above LevelError.
// Arguments are handled in the manner of fmt.Printf.
func (log *Log) Errorf(format string, value ...interface{}) {
	log.logf(LevelError, format, value...)
}

// Warning logs record if given logger level is equal or above LevelWarning.
// Arguments are handled in the manner of fmt.Print.
func (log *Log) Warning(value ...interface{}) {
	log.log(LevelWarning, value...)
}

// Warningf logs record if given logger level is equal or above LevelWarning.
// Arguments are handled in the manner of fmt.Printf.
func (log *Log) Warningf(format string, value ...interface{}) {
	log.logf(LevelWarning, format, value...)
}

// Print is pseudonym for Info
func (log *Log) Print(value ...interface{}) {
	log.log(LevelInfo, value...)
}

// Printf is pseudonym for Infof
func (log *Log) Printf(format string, value ...interface{}) {
	log.logf(LevelInfo, format, value...)
}

// Info logs record if given logger level is equal or above LevelInfo.
// Arguments are handled in the manner of fmt.Print.
func (log *Log) Info(value ...interface{}) {
	log.log(LevelInfo, value...)
}

// Infof logs record if given logger level is equal or above LevelInfo.
// Arguments are handled in the manner of fmt.Printf.
func (log *Log) Infof(format string, value ...interface{}) {
	log.logf(LevelInfo, format, value...)
}

// Debug logs record if given logger level is equal or above LevelDebug.
// Arguments are handled in the manner of fmt.Print.
func (log *Log) Debug(value ...interface{}) {
	log.log(LevelDebug, value...)
}

// Debugf logs record if given logger level is equal or above LevelDebug.
// Arguments are handled in the manner of fmt.Printf.
func (log *Log) Debugf(format string, value ...interface{}) {
	log.logf(LevelDebug, format, value...)
}

// Trace logs record if given logger level is equal or above LevelTrace.
// Arguments are handled in the manner of fmt.Print.
func (log *Log) Trace(value ...interface{}) {
	log.log(LevelTrace, value...)
}

// Tracef logs record if given logger level is equal or above LevelTrace.
// Arguments are handled in the manner of fmt.Printf.
func (log *Log) Tracef(format string, value ...interface{}) {
	log.logf(LevelTrace, format, value...)
}

// SetPrefix of given logger, prefix placeholder should be used in logger
// format.
func (log *Log) SetPrefix(prefix string) {
	log.prefix = prefix
}

// NewChild of given logger, child inherit level, format and output options.
func (log *Log) NewChild() *Log {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	child := NewLog()
	child.SetOutput(log.output)
	child.SetLevel(log.level)
	child.SetFormat(log.format)

	log.children = append(log.children, child)

	return child
}

// NewChildWithPrefix of given logger, child inherit level, format and output
// options.
func (log *Log) NewChildWithPrefix(prefix string) *Log {
	child := log.NewChild()
	child.SetPrefix(prefix)
	return child
}
