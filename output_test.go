package lorg

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutput_SetLevelWriterCondition_SetsAdditionalWriter(t *testing.T) {
	test := assert.New(t)

	var buffer0 bytes.Buffer
	var buffer1 bytes.Buffer
	var buffer2 bytes.Buffer

	logger := NewLog()
	logger.SetOutput(
		NewOutput(&buffer0).(*output).SetLevelWriterCondition(
			LevelWarning, &buffer1, &buffer2,
		),
	)
	logger.SetFormat(
		NewFormat("${level} %s"),
	)

	logger.Info("1")
	logger.Warning("2")

	test.EqualValues(buffer0.String(), "INFO 1\n")
	test.EqualValues(buffer1.String(), "WARNING 2\n")
	test.EqualValues(buffer2.String(), "WARNING 2\n")
}
