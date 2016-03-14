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
	}

	return "UNKNOWN"
}
