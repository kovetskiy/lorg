package shortcuts

import (
	"github.com/kovetskiy/lorg"
)

var (
	logger lorg.Logger = lorg.NewLog()
)

func SetLogger(log lorg.Logger) {
	logger = log
}

func Fatalf(format string, values ...interface{}) {
	logger.Fatalf(format, values...)
}

func Errorf(format string, values ...interface{}) {
	logger.Errorf(format, values...)
}

func Warningf(format string, values ...interface{}) {
	logger.Warningf(format, values...)
}

func Infof(format string, values ...interface{}) {
	logger.Infof(format, values...)
}

func Debugf(format string, values ...interface{}) {
	logger.Debugf(format, values...)
}

func Tracef(format string, values ...interface{}) {
	logger.Tracef(format, values...)
}

func Fatal(values ...interface{}) {
	logger.Fatal(values...)
}

func Error(values ...interface{}) {
	logger.Error(values...)
}

func Warning(values ...interface{}) {
	logger.Warning(values...)
}

func Info(values ...interface{}) {
	logger.Info(values...)
}

func Debug(values ...interface{}) {
	logger.Debug(values...)
}

func Trace(values ...interface{}) {
	logger.Trace(values...)
}
