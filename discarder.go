package lorg

// NewDiscarder returns new Discarder instance which implements Logger
// interface but have one important feature, Discarder actually do nothing and
// doesn't log anything.
//
// It's very useful in packages, which want to have a opportunity to log debug
// messages, but by default should not log anything.
func NewDiscarder() Logger {
	return new(discarder)
}

// ensure that discarder implements Logger interface.
var _ Logger = (*discarder)(nil)

type discarder struct{}

func (*discarder) Fatal(_ ...interface{})              {}
func (*discarder) Fatalf(_ string, _ ...interface{})   {}
func (*discarder) Error(_ ...interface{})              {}
func (*discarder) Errorf(_ string, _ ...interface{})   {}
func (*discarder) Warning(_ ...interface{})            {}
func (*discarder) Warningf(_ string, _ ...interface{}) {}
func (*discarder) Print(_ ...interface{})              {}
func (*discarder) Printf(_ string, _ ...interface{})   {}
func (*discarder) Info(_ ...interface{})               {}
func (*discarder) Infof(_ string, _ ...interface{})    {}
func (*discarder) Debug(_ ...interface{})              {}
func (*discarder) Debugf(_ string, _ ...interface{})   {}
