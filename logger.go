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
	Fatal(values ...interface{})
	Fatalf(format string, values ...interface{})

	Error(values ...interface{})
	Errorf(format string, values ...interface{})

	Warning(values ...interface{})
	Warningf(format string, values ...interface{})

	Print(values ...interface{})
	Printf(format string, values ...interface{})

	Info(values ...interface{})
	Infof(format string, values ...interface{})

	Debug(values ...interface{})
	Debugf(format string, values ...interface{})

	Trace(values ...interface{})
	Tracef(format string, values ...interface{})
}
