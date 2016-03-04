package lorg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFormat_ReturnsFormatWithDefaultFiends(t *testing.T) {
	format := NewFormat(``)

	assert.Len(t, format.placeholders, len(defaultPlaceholders))

	keys := []string{}
	for key := range defaultPlaceholders {
		keys = append(keys, key)
	}

	for _, key := range keys {
		assert.Contains(t, format.placeholders, key)
	}
}

func TestFormat_ImplementsFormatterInterface(t *testing.T) {
	assert.Implements(t, (*Formatter)(nil), &Format{})
}

func TestFormat_GetPlaceholders_ReturnsPlaceholdersField(t *testing.T) {
	rawFormat := `${place_foo} ${place_bar:arg} ${unknown}`
	format := NewFormat(rawFormat)

	placeholders := map[string]Placeholder{
		"place_foo": func(_ Level, _ string) string { return "holder foo" },
		"place_bar": func(_ Level, _ string) string { return "holder bar" },
	}

	format.placeholders = placeholders

	assert.Equal(t, placeholders, format.GetPlaceholders())
}

func TestFormat_SetPlaceholders_ChangesPlaceholdersField(t *testing.T) {
	// https://golang.org/doc/devel/weekly.html#2011-11-18
	//
	// Map and function value comparisons are now disallowed (except for
	// comparison with nil) as per the Go 1 plan. Function equality was
	// problematic in some contexts and map equality compares pointers, not the
	// maps' content.
	//
	// Okay, go can't compare functions, so we can't test that SetPlaceholder
	// sets exactly specified Placeholder into placeholders[key], we can check
	// only key existing.
	format := NewFormat(``)

	placeholders := map[string]Placeholder{
		"place_foo": func(_ Level, _ string) string { return "foo" },
	}

	format.SetPlaceholders(placeholders)
	assert.Contains(t, format.placeholders, "place_foo")

	anotherPlaceholders := map[string]Placeholder{
		"place_bar": func(_ Level, _ string) string { return "bar" },
	}

	format.SetPlaceholders(anotherPlaceholders)
	assert.NotContains(t, format.placeholders, "place_foo")
	assert.Contains(t, format.placeholders, "place_bar")
}

func TestFormat_SetPlaceholder_ChangesPlaceholdersField(t *testing.T) {
	// please see comment about testing SetPlaceholders function:
	// TestFormat_SetPlaceholders_ChangesPlaceholdersField
	format := NewFormat(``)

	placeholderFoo := Placeholder(
		func(_ Level, _ string) string { return "foo" },
	)
	format.SetPlaceholder("place_foo", placeholderFoo)

	assert.Equal(t, len(defaultPlaceholders)+1, len(format.placeholders))
	assert.Contains(t, format.placeholders, "place_foo")

	placeholderBar := Placeholder(
		func(_ Level, _ string) string { return "bar" },
	)
	format.SetPlaceholder("place_bar", placeholderBar)

	assert.Equal(t, len(defaultPlaceholders)+2, len(format.placeholders))
	assert.Contains(t, format.placeholders, "place_foo")
	assert.Contains(t, format.placeholders, "place_bar")
}

func TestFormat_Render_UsesPlaceholdersFieldForMatchingPlaceholders(
	t *testing.T,
) {
	rawFormat := `plain text ${place_foo} foo ${place_bar:barvalue}`
	format := NewFormat(rawFormat)

	placeholders := map[string]Placeholder{
		"place_foo": func(_ Level, _ string) string { return "foo" },
	}

	format.SetPlaceholders(placeholders)
	format.Render(LevelWarning)

	assert.Equal(
		t,
		[]string{`${place_foo}`},
		getReplacementsValues(format.replacements),
	)

	format.SetPlaceholders(
		map[string]Placeholder{
			"place_foo": func(_ Level, _ string) string { return "foo" },
			"place_bar": func(_ Level, _ string) string { return "bar" },
		},
	)

	format.Reset()

	format.Render(LevelWarning)

	assert.Equal(
		t,
		[]string{`${place_foo}`, `${place_bar:barvalue}`},
		getReplacementsValues(format.replacements),
	)
}

func TestFormat_Render_PlaceholderRegexpMatching(t *testing.T) {
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

func TestFormat_Reset_DesolatesReplacementAndUnsetsCompiledFlag(
	t *testing.T,
) {
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
