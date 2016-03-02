package lorg

// Logger is the interface which implemented by Log structure
//
// It's very usable if you creates package which want to have opportunity to
// log debug data.
// In a package you should create unexported global variable like as following:
//
//     var log lorg.Logger = lorg.Discarder
//
// and exported function for setting a logger:
//
//     func SetLog(logger lorg.Logger) {
//        log = logger
//     }
//
// After this you can use `log.Debug` and `log.Debugf` in your package code,
// and since the log variable is lorg.Discarder by default your package will
// not log anything until package-user want to see package debug messages and
// sets logger using yourPackage.SetLog(someLogger)
type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Warning(args ...interface{})
	Warningf(format string, args ...interface{})

	Print(args ...interface{})
	Printf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}
