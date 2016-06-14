package lorg

// Level describes all available log levels for log records.
type Level int

const (
	// LevelFatal should be used for messages with FATAL severity.
	LevelFatal Level = iota

	// LevelError should be used for messages with ERROR severity.
	LevelError

	// LevelWarning should be used for messages with WARNING severity.
	LevelWarning

	// LevelInfo should be used for messages with INFO severity.
	LevelInfo

	// LevelDebug should be used for messages with DEBUG severity.
	LevelDebug

	// LevelTrace should be used for messages with TRACE severity.
	LevelTrace
)

// String returns the string representation of a given logging level.
func (level Level) String() string {
	switch level {
	case LevelFatal:
		return "FATAL"
	case LevelError:
		return "ERROR"
	case LevelWarning:
		return "WARNING"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	case LevelTrace:
		return "TRACE"
	}

	return "UNKNOWN"
}

// String returns the short string representation of a given logging level.
func (level Level) StringShort() string {
	switch level {
	case LevelWarning:
		return "WARN"
	}

	return level.String()
}
