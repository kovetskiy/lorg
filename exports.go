package lorg

import (
	"io"
)

var (
	logger = NewLog()
)

// SetLevel sets the logging level for the given log.
// After setting level, logger will log records with same level or above.
//
// Running SetLevel it's not required operation, by default Log instance
// creates with INFO level, so levels above (warn, err, fatal, info) will be
// logged also.
func SetLevel(level Level) {
	logger.SetLevel(level)
}

// SetFormat sets the logging format for the given log.
// All log records will be formatted using specified formatter.
//
// Running SetFormat it's not required operation, by default Log instance
// creates with default Format instance, which have format with level and date
// placeholders.
// See: DefaultFormatting
func SetFormat(format Formatter) {
	logger.SetFormat(format)
}

// SetOutput sets output of given log instance, the log records are
// fmt.Printfd'd to specified io.Writer.
//
// Running SetOutput it's not required operation, by default Log instance
// logs all records to stderr (os.Stderr)
func SetOutput(output io.Writer) {
	logger.SetOutput(output)
}

// Fatal logs record if given logger level is equal or above LevelFatal, and
// calls os.Exit(1) after logging.
// Arguments are handled in the manner of fmt.Print.
func Fatal(value ...interface{}) {
	logger.log(LevelFatal, value...)
	Exiter(1)
}

// Fatalf logs record if given logger level is equal or above LevelFatal, and
// calls os.Exit(1) after logging.
// Arguments are handled in the manner of fmt.Print.
func Fatalf(format string, value ...interface{}) {
	logger.logf(LevelFatal, format, value...)
	Exiter(1)
}

// Error logs record if given logger level is equal or above LevelError.
// Arguments are handled in the manner of fmt.Print.
func Error(value ...interface{}) {
	logger.log(LevelError, value...)
}

// Errorf logs record if given logger level is equal or above LevelError.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, value ...interface{}) {
	logger.logf(LevelError, format, value...)
}

// Warning logs record if given logger level is equal or above LevelWarning.
// Arguments are handled in the manner of fmt.Print.
func Warning(value ...interface{}) {
	logger.log(LevelWarning, value...)
}

// Warningf logs record if given logger level is equal or above LevelWarning.
// Arguments are handled in the manner of fmt.Printf.
func Warningf(format string, value ...interface{}) {
	logger.logf(LevelWarning, format, value...)
}

// Print is pseudonym for Info
func Print(value ...interface{}) {
	logger.log(LevelInfo, value...)
}

// Printf is pseudonym for Infof
func Printf(format string, value ...interface{}) {
	logger.logf(LevelInfo, format, value...)
}

// Info logs record if given logger level is equal or above LevelInfo.
// Arguments are handled in the manner of fmt.Print.
func Info(value ...interface{}) {
	logger.log(LevelInfo, value...)
}

// Infof logs record if given logger level is equal or above LevelInfo.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, value ...interface{}) {
	logger.logf(LevelInfo, format, value...)
}

// Debug logs record if given logger level is equal or above LevelDebug.
// Arguments are handled in the manner of fmt.Print.
func Debug(value ...interface{}) {
	logger.log(LevelDebug, value...)
}

// Debugf logs record if given logger level is equal or above LevelDebug.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, value ...interface{}) {
	logger.logf(LevelDebug, format, value...)
}

// Trace logs record if given logger level is equal or above LevelTrace.
// Arguments are handled in the manner of fmt.Print.
func Trace(value ...interface{}) {
	logger.log(LevelTrace, value...)
}

// Tracef logs record if given logger level is equal or above LevelTrace.
// Arguments are handled in the manner of fmt.Printf.
func Tracef(format string, value ...interface{}) {
	logger.logf(LevelTrace, format, value...)
}
