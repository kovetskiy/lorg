package lorg

type Level int

const (
	LevelFatal Level = iota
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
)

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
