package lorg

var Discarder Logger = new(discarder)

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
