package lorg

import "strings"

// Format is the actual Formatter which used by Log structure for formatting
// log records before writing log records into Log.output.
//
// Do not instantiate Format instance without using NewFormat.
type Format struct {
	rawFormat    string
	compiled     bool
	replacements []replacement
	placeholders map[string]Placeholder
}

type replacement struct {
	value          string
	placeholder    Placeholder
	placeholderArg string
}

// NewFormat creates Format instance with specified formatting and default
// placeholders: level (PlaceholderLevel), date (PlaceholderDate), line
// (PlaceholderLine) and file (PlaceholderFile).
//
// Format placeholders can be changed or added using SetPlaceholders or SetPlaceholder
// methods.
func NewFormat(rawFormat string) *Format {
	format := &Format{
		rawFormat: rawFormat,
	}

	// we are should not assing format.placeholders to defaultPlaceholders
	// because maps in go passes by reference.

	format.SetPlaceholders(defaultPlaceholders)

	return format
}

// SetPlaceholder sets specified placeholder with specified placeholder name
// for given format.
func (format *Format) SetPlaceholder(name string, placeholder Placeholder) {
	format.placeholders[name] = placeholder
}

// SetPlaceholders sets specified placeholders for given format.
func (format *Format) SetPlaceholders(placeholders map[string]Placeholder) {
	format.placeholders = map[string]Placeholder{}

	for placeholderName, placeholder := range placeholders {
		format.placeholders[placeholderName] = placeholder
	}
}

// GetPlaceholders returns placeholders of given format.
func (format *Format) GetPlaceholders() map[string]Placeholder {
	return format.placeholders
}

// Reset resets state of given format.
func (format *Format) Reset() {
	format.replacements = []replacement{}
	format.compiled = false
}

// Render generates string which will be used by Log instance.
// Here is logLevel property just for a placeholders which want to show
// logging level, logLevel will be passed to all ran placeholders.
func (format *Format) Render(logLevel Level) string {
	if !format.compiled {
		format.compile()
	}

	rendered := format.rawFormat
	for _, replacement := range format.replacements {
		placeholderValue := replacement.placeholder(
			logLevel,
			replacement.placeholderArg,
		)

		rendered = strings.Replace(
			rendered, replacement.value, placeholderValue, 1,
		)
	}

	return rendered
}

func (format *Format) compile() {
	// reset compiled replacements
	format.Reset()

	matches := rePlaceholder.FindAllStringSubmatch(format.rawFormat, -1)
	for _, match := range matches {
		var (
			replacementValue = match[0]
			placeholderName  = match[1]
			placeholderArg   = match[3]
		)

		placeholder, ok := format.placeholders[placeholderName]
		if !ok {
			// placeholder with specified name not found
			continue
		}

		newReplacement := replacement{
			value:          replacementValue,
			placeholder:    placeholder,
			placeholderArg: placeholderArg,
		}

		format.replacements = append(format.replacements, newReplacement)
	}

	format.compiled = true
}
