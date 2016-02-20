package lorg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatImplementsFormatter(t *testing.T) {
	assert.Implements(t, (*Formatter)(nil), NewFormat(``))
}

func TestFormatDefaultPlaceholders(t *testing.T) {
	format := NewFormat(``)
	assert.Equal(t, defaultPlaceholders, format.placeholders)
}

func TestFormatGetPlaceholders(t *testing.T) {
	rawFormat := `${place_foo} ${place_bar:arg} ${unknown}`
	format := NewFormat(rawFormat)

	placeholders := map[string]Placeholder{
		"place_foo": func(_ Level, _ string) string { return "holder foo" },
		"place_bar": func(_ Level, _ string) string { return "holder bar" },
	}

	format.SetPlaceholders(placeholders)

	assert.Equal(t, placeholders, format.GetPlaceholders())
	assert.Equal(t, placeholders, format.placeholders)
}

func TestFormatPlaceholdersMatching(t *testing.T) {
	type testcase struct {
		format               string
		expectedReplacements []string
	}

	testcases := []testcase{
		{
			`text`,
			[]string{},
		},
		{
			`${place_foo}`,
			[]string{
				`${place_foo}`,
			},
		},
		{
			`before${place_foo}after`,
			[]string{
				`${place_foo}`,
			},
		},
		{
			`before space ${place_foo} space after`,
			[]string{
				`${place_foo}`,
			},
		},
		{
			`before${place_foo}`,
			[]string{
				`${place_foo}`,
			},
		},
		{
			`before space ${place_foo}`,
			[]string{
				`${place_foo}`,
			},
		},
		{
			`${place_foo}after`,
			[]string{
				`${place_foo}`,
			},
		},
		{
			`${place_foo} space after`,
			[]string{
				`${place_foo}`,
			},
		},
		{
			`${place_foo} ${unknown}`,
			[]string{
				`${place_foo}`,
			},
		},
		{
			`${place_foo} ${place_foo:1}`,
			[]string{
				`${place_foo}`, `${place_foo:1}`,
			},
		},
		{
			`${place_foo} ${place_foo:1:2}`,
			[]string{
				`${place_foo}`, `${place_foo:1:2}`,
			},
		},
		{
			`${place_foo} ${unknown} ${place_foo:1}`,
			[]string{
				`${place_foo}`, `${place_foo:1}`,
			},
		},
		{
			`${unknown}${place_foo} ${place_foo:1}after`,
			[]string{
				`${place_foo}`, `${place_foo:1}`,
			},
		},
		{
			`${place_foo} ${place_foo:1} ${place_foo}`,
			[]string{
				`${place_foo}`, `${place_foo:1}`, `${place_foo}`,
			},
		},
		{
			`${place_foo} ${place_foo:1} ${place_foo:1:2:3:4}`,
			[]string{
				`${place_foo}`, `${place_foo:1}`, `${place_foo:1:2:3:4}`,
			},
		},
		{
			`${place_foo} ${place_foo:1} ${place_bar}}`,
			[]string{
				`${place_foo}`, `${place_foo:1}`, `${place_bar}`,
			},
		},
	}

	placeholders := map[string]Placeholder{
		"place_foo": func(_ Level, _ string) string { return "holder foo" },
		"place_bar": func(_ Level, _ string) string { return "holder bar" },
	}

	for _, testcase := range testcases {
		format := NewFormat(testcase.format)
		format.SetPlaceholders(placeholders)
		format.Render(LevelError)

		assert.Equal(
			t,
			testcase.expectedReplacements,
			getReplacementsValues(format.replacements),
			"format: %s", testcase.format,
		)
	}
}

func TestFormatSetPlaceholders(t *testing.T) {
	rawFormat := `plain text ${place_foo} foo ${place_bar:barvalue}`
	format := NewFormat(rawFormat)

	format.SetPlaceholders(
		map[string]Placeholder{
			"place_foo": func(_ Level, _ string) string { return "foo" },
		},
	)

	format.Render(LevelWarning)

	assert.Equal(
		t,
		[]string{`${place_foo}`},
		getReplacementsValues(format.replacements),
	)

	format.SetPlaceholder(
		"place_bar", func(_ Level, _ string) string { return "bar" },
	)

	format.Reset()

	format.Render(LevelWarning)

	assert.Equal(
		t,
		[]string{`${place_foo}`, `${place_bar:barvalue}`},
		getReplacementsValues(format.replacements),
	)
}

func TestFormatRenderAndReset(t *testing.T) {
	rawFormat := `plain text ${place_foo} foo ${place_bar:barvalue}`
	format := NewFormat(rawFormat)

	format.SetPlaceholders(
		map[string]Placeholder{
			"place_foo": func(_ Level, _ string) string { return "foo" },
		},
	)

	format.Render(LevelWarning)

	assert.NotEmpty(t, format.replacements)
	assert.True(t, format.compiled)

	format.Reset()

	assert.Empty(t, format.replacements)
	assert.False(t, format.compiled)
}
