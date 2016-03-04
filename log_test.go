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
	log := NewLog()

	assert.Equal(
		t, defaultLevel, log.level,
		"Log object created with wrong default logging level",
	)

	assert.Equal(
		t, defaultFormat, log.format,
		"Log object created with wrong default logging format",
	)

	assert.Equal(
		t, defaultOutput, log.output,
		"Log object created with wrong default logging output",
	)
}

func TestLog_ImplementsLoggerInterface(t *testing.T) {
	assert.Implements(t, (*Logger)(nil), &Log{})
}

func TestLog_SetFormat_ChangesFormatField(t *testing.T) {
	mock := new(mockFormat)

	log := NewLog()
	log.SetFormat(mock)

	assert.Equal(t, mock, log.format)
}

func TestLog_SetLevel_ChangesLevelField(t *testing.T) {
	log := NewLog()
	log.SetLevel(LevelWarning)

	assert.Equal(t, LevelWarning, log.level)
}

func TestLog_SetOutput_ChangesOutputField(t *testing.T) {
	log := NewLog()
	log.SetOutput(ioutil.Discard)

	assert.Equal(t, ioutil.Discard, log.output)
}

func TestLog_LoggingFunctions_CallsFormatRender(t *testing.T) {
	mock := new(mockFormat)

	log := NewLog()
	log.SetFormat(mock)

	log.Print(getFilenameAndLine())

	assert.Equal(t, 1, mock.callsRender)

	log.Printf(getFilenameAndLine())
	log.Error(getFilenameAndLine())
	log.Errorf(getFilenameAndLine())
	log.Warning(getFilenameAndLine())
	log.Warningf(getFilenameAndLine())
	log.Info(getFilenameAndLine())
	log.Infof(getFilenameAndLine())

	assert.Equal(t, 8, mock.callsRender)

	log.SetLevel(LevelDebug)

	log.Debug(getFilenameAndLine())
	log.Debugf(getFilenameAndLine())

	assert.Equal(t, 10, mock.callsRender)
}

func TestLog_LoggingFunctions_LogsRecordsWithSameLevelOrAbove(
	t *testing.T,
) {
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

		assert.Equal(
			t, (index+1)*2, mock.callsRender,
			"level: %s", setting.level,
		)
	}
}
