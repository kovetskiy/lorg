package lorg

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockFormat struct {
	lastLogLevel Level
	callsRender  int
}

func (mock *mockFormat) SetPlaceholders(_ map[string]Placeholder) {
	panic("should not be called")
}

func (mock *mockFormat) SetPlaceholder(_ string, _ Placeholder) {
	panic("should not be called")
}

func (mock *mockFormat) GetPlaceholders() map[string]Placeholder {
	panic("should not be called")
}

func (mock *mockFormat) SetPrefix(_ string) {
	panic("should be not called")
}

func (mock *mockFormat) Render(logLevel Level, prefix string) string {
	mock.callsRender++
	mock.lastLogLevel = logLevel
	return "[testcase] %s"
}

func (mock *mockFormat) Reset() {}

func TestNewLog_ReturnsLogWithDefaultFields(t *testing.T) {
	test := assert.New(t)

	log := NewLog()

	test.Equal(
		defaultLevel, log.level,
		"Log object created with wrong default logging level",
	)

	test.Equal(
		defaultFormat, log.format,
		"Log object created with wrong default logging format",
	)

	test.Equal(
		defaultOutput, log.output,
		"Log object created with wrong default logging output",
	)
}

func TestLog_ImplementsLoggerInterface(t *testing.T) {
	test := assert.New(t)

	test.Implements((*Logger)(nil), &Log{})
}

func TestLog_SetFormat_ChangesFormatField(t *testing.T) {
	test := assert.New(t)

	mock := new(mockFormat)

	log := NewLog()
	log.SetFormat(mock)

	test.Equal(mock, log.format)
}

func TestLog_SetLevel_ChangesLevelField(t *testing.T) {
	test := assert.New(t)

	log := NewLog()
	log.SetLevel(LevelWarning)

	test.Equal(LevelWarning, log.level)
}

func TestLog_SetOutput_ChangesOutputField(t *testing.T) {
	test := assert.New(t)

	log := NewLog()
	log.SetOutput(ioutil.Discard)

	test.Equal(ioutil.Discard, log.output.(*Output).conditions[LevelDebug][0])
}

func TestLog_LoggingFunctions_CallsFormatRender(t *testing.T) {
	test := assert.New(t)

	mock := new(mockFormat)

	log := NewLog()
	log.SetFormat(mock)

	log.Print(getFilenameAndLine())

	test.Equal(1, mock.callsRender)

	log.Printf(getFilenameAndLine())
	log.Error(getFilenameAndLine())
	log.Errorf(getFilenameAndLine())
	log.Warning(getFilenameAndLine())
	log.Warningf(getFilenameAndLine())
	log.Info(getFilenameAndLine())
	log.Infof(getFilenameAndLine())

	test.Equal(8, mock.callsRender)

	log.SetLevel(LevelDebug)

	log.Debug(getFilenameAndLine())
	log.Debugf(getFilenameAndLine())

	test.Equal(10, mock.callsRender)
}

func TestLog_LoggingFunctions_LogsRecordsWithSameLevelOrAbove(
	t *testing.T,
) {
	test := assert.New(t)

	log := NewLog()

	// Fatal tested in special function
	levels := []struct {
		level Level
		log   func(...interface{})
		logf  func(string, ...interface{})
	}{
		{LevelError, log.Error, log.Errorf},
		{LevelWarning, log.Warning, log.Warningf},
		{LevelInfo, log.Info, log.Infof},
		{LevelDebug, log.Debug, log.Debugf},
	}

	for index, setting := range levels {
		mock := new(mockFormat)
		log.SetFormat(mock)
		log.SetLevel(setting.level)

		for _, checking := range levels {
			checking.log(
				setting.level.String(),
				checking.level.String(),
				getFilenameAndLine(),
			)

			checking.logf(
				"%s %s %s",
				setting.level,
				checking.level,
				getFilenameAndLine(),
			)
		}

		test.Equal(
			(index+1)*2, mock.callsRender,
			"level: %s", setting.level,
		)
	}
}

func TestLog_NewChild_InheritsLevelValue(t *testing.T) {
	test := assert.New(t)

	log := NewLog()
	log.SetLevel(LevelDebug)

	child := log.NewChild()
	test.Equal(child.level, log.level)
}

func TestLog_NewChild_InheritsOutputValue(t *testing.T) {
	test := assert.New(t)

	log := NewLog()
	child := log.NewChild()

	test.Equal(address(child.output), address(log.output))
}

func TestLog_SetLevel_ChangesChildrenLevelField(t *testing.T) {
	test := assert.New(t)

	log := NewLog()
	child := log.NewChild()
	child2 := log.NewChild()

	log.SetLevel(LevelTrace)

	test.Equal(child.level, log.level)
	test.Equal(child2.level, log.level)
}

func TestLog_NewChild_ChildCantChangeParentLevel(t *testing.T) {
	test := assert.New(t)

	log := NewLog()
	child := log.NewChild()

	log.SetLevel(LevelTrace)
	child.SetLevel(LevelDebug)

	test.NotEqual(child.level, log.level)
}

func TestLog_NewChildWithPrefix_ReturnsLoggerWithPrefix(t *testing.T) {
	test := assert.New(t)

	var buffer bytes.Buffer

	log := NewLog()
	log.SetOutput(&buffer)
	log.SetFormat(NewFormat(`${prefix}%s`))
	child := log.NewChildWithPrefix("child")
	subchild := child.NewChildWithPrefix("subchild")

	log.Info("1")
	child.Info("2")
	subchild.Info("3")
	child.Info("4")
	log.Info("5")

	test.Equal("1\nchild 2\nsubchild 3\nchild 4\n5\n", buffer.String())
}

func TestLog_IndentLines(t *testing.T) {
	test := assert.New(t)

	var buffer bytes.Buffer

	log := NewLog()
	log.SetOutput(&buffer)
	log.SetFormat(NewFormat(`[blah] %s`))
	log.SetIndentLines(true)

	log.Info("1")
	log.Info("2\n3\n4")

	test.Equal(
		`[blah] 1
[blah] 2
       3
       4
`, buffer.String())
}

func TestLog_IndentLinesWithTrimStyle(t *testing.T) {
	test := assert.New(t)

	var buffer bytes.Buffer

	log := NewLog()
	log.SetOutput(&buffer)
	log.SetFormat(
		NewFormat("\x1b[48;5;2mblah: %s"),
	)

	log.SetIndentLines(true)

	log.Info("1")
	log.Info("2\n3\n4")

	test.Equal(
		"\x1b[48;5;2m"+
			"blah: 1\n"+
			"\x1b[48;5;2m"+
			"blah: 2\n"+
			"      3\n"+
			"      4\n",
		buffer.String(),
	)
}

func address(target interface{}) uintptr {
	value := reflect.ValueOf(target)
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		value = value.Elem()
	}

	return value.UnsafeAddr()
}
