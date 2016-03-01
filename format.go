package lorg

import "strings"

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

func NewFormat(rawFormat string) *Format {
	format := &Format{
		rawFormat:    rawFormat,
		placeholders: defaultPlaceholders,
	}

	return format
}

func (format *Format) SetPlaceholder(name string, placeholder Placeholder) {
	format.placeholders[name] = placeholder
}

func (format *Format) SetPlaceholders(placeholders map[string]Placeholder) {
	format.placeholders = placeholders
}

func (format *Format) GetPlaceholders() map[string]Placeholder {
	return format.placeholders
}

func (format *Format) Reset() {
	format.replacements = []replacement{}
	format.compiled = false
}

// here is logLevel property just for a placeholders which want to show
// logging level
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
