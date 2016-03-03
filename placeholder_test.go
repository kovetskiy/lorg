package lorg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// placeholderFabric can fabricate a function which has same signature
// as Placeholder, but fabricated function logs entry about running
// placeholder to fabric's log and returns same entry.
type placeholderFabric struct {
	log []string
}

func (fabric *placeholderFabric) fabricate(name string) Placeholder {
	// create a function and pass placeholder name by reference
	return func(_ Level, value string) string {
		message := fmt.Sprintf("[%s@%s]", name, value)
		fabric.log = append(fabric.log, message)
		return message
	}
}

func TestPlaceholderRunning(t *testing.T) {
	fabric := new(placeholderFabric)

	format := NewFormat(
		`fmt: ${place_foo} ${place_foo:1} ${place_bar:a b c:d e f}`,
	)

	format.SetPlaceholders(
		map[string]Placeholder{
			"place_foo": fabric.fabricate("place_foo"),
			"place_bar": fabric.fabricate("place_bar"),
		},
	)

	rendered := format.Render(LevelWarning)

	assert.Equal(
		t,
		"fmt: [place_foo@] [place_foo@1] [place_bar@a b c:d e f]",
		rendered,
	)

	expectedFabricLog := []string{
		"[place_foo@]",
		"[place_foo@1]",
		"[place_bar@a b c:d e f]",
	}

	assert.Equal(
		t, expectedFabricLog, fabric.log,
		"Format runs placeholders in wrong order",
	)

	// now let's fake placeholder for place_bar and check that placeholder
	// was launched, and old place_bar placeholder wasn't launched
	format.SetPlaceholder("place_bar", fabric.fabricate("fakebar"))
	format.Reset()

	rendered = format.Render(LevelWarning)

	assert.Equal(
		t,
		"fmt: [place_foo@] [place_foo@1] [fakebar@a b c:d e f]",
		rendered,
	)

	expectedFabricLog = append(
		expectedFabricLog,
		"[place_foo@]",
		"[place_foo@1]",
		"[fakebar@a b c:d e f]",
	)

	assert.Equal(
		t, expectedFabricLog, fabric.log,
		"Format runs placeholders in wrong order",
	)
}

func TestFormatPassesLogLevelToPlaceholders(t *testing.T) {
	format := NewFormat(`${place_foo}`)

	var placeholderLogLevel Level = -1

	format.SetPlaceholder(
		"place_foo",
		func(logLevel Level, _ string) string {
			placeholderLogLevel = logLevel
			return ""
		},
	)

	format.Render(LevelDebug)

	assert.Equal(
		t, LevelDebug, placeholderLogLevel,
		"log level doesn't passed to placeholder",
	)
}
