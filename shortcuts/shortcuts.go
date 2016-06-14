package shortcuts

import (
	"github.com/kovetskiy/lorg"
)

var (
	Logger = lorg.NewLog()
)

func Fatalf(format string, values ...interface{}) {
	Logger.Fatalf(format, values...)
}

func Errorf(format string, values ...interface{}) {
	Logger.Errorf(format, values...)
}

func Warningf(format string, values ...interface{}) {
	Logger.Warningf(format, values...)
}

func Infof(format string, values ...interface{}) {
	Logger.Infof(format, values...)
}

func Debugf(format string, values ...interface{}) {
	Logger.Debugf(format, values...)
}

func Tracef(format string, values ...interface{}) {
	Logger.Tracef(format, values...)
}
