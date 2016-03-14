package lorg

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockFormat struct {
	lastLogLevel Level
	callsRender  int
}

func (mock *mockFormat) SetPlaceholders(_ map[string]Placeholder) {
	panic("SetPlaceholders method should not be called")
}

func (mock *mockFormat) SetPlaceholder(_ string, _ Placeholder) {
	panic("SetPlaceholder method should not be called")
}

func (mock *mockFormat) GetPlaceholders() map[string]Placeholder {
	panic("GetPlaceholder method should not be called")
}

func (mock *mockFormat) Render(logLevel Level) string {
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

	test.Equal(ioutil.Discard, log.output)
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
