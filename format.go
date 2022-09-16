package lorg

import (
	"strings"
	"sync"
)

// Format is the actual Formatter which used by Log structure for formatting
// log records before writing log records into Log.output.
//
// Do not instantiate Format instance without using NewFormat.
type Format struct {
	formatting       string
	compiled         bool
	compileMutex     *sync.Mutex
	replacements     []replacement
	placeholders     map[string]Placeholder
	placeholderMutex *sync.RWMutex
}

type replacement struct {
	value            string
	placeholder      Placeholder
	placeholderValue string
}

// NewFormat creates Format instance with specified formatting and default
// placeholders: level (PlaceholderLevel), date (PlaceholderDate), line
// (PlaceholderLine) and file (PlaceholderFile).
//
// Format placeholders can be changed or added using SetPlaceholders or
// SetPlaceholder methods.
func NewFormat(formatting string) *Format {
	format := &Format{
		formatting:       formatting,
		placeholderMutex: &sync.RWMutex{},
		compileMutex:     &sync.Mutex{},
	}

	// we are should not assing format.placeholders to defaultPlaceholders
	// because maps in go passes by reference.

	format.SetPlaceholders(DefaultPlaceholders)

	return format
}

// SetPlaceholder sets specified placeholder with specified placeholder name
// for given format.
func (format *Format) SetPlaceholder(name string, placeholder Placeholder) {
	format.Reset()

	format.placeholderMutex.Lock()
	format.placeholders[name] = placeholder
	format.placeholderMutex.Unlock()
}

// SetPlaceholders sets specified placeholders for given format.
func (format *Format) SetPlaceholders(placeholders map[string]Placeholder) {
	format.Reset()

	format.placeholderMutex.Lock()

	format.placeholders = map[string]Placeholder{}
	for placeholderName, placeholder := range placeholders {
		format.placeholders[placeholderName] = placeholder
	}

	format.placeholderMutex.Unlock()
}

// GetPlaceholders returns placeholders of given format.
func (format *Format) GetPlaceholders() map[string]Placeholder {
	return format.placeholders
}

// Reset resets state of given format.
func (format *Format) Reset() {
	format.placeholderMutex.Lock()

	format.replacements = []replacement{}
	format.compiled = false
	cache.reset()

	format.placeholderMutex.Unlock()
}

// Render generates string which will be used by Log instance.
// Here is logLevel property just for a placeholders which want to show
// logging level, logLevel will be passed to all ran placeholders.
func (format *Format) Render(logLevel Level, prefix string) string {
	format.compileMutex.Lock()
	if !format.compiled {
		format.compile()
	}
	format.compileMutex.Unlock()

	format.placeholderMutex.RLock()
	rendered := format.formatting
	for _, replacement := range format.replacements {
		var placeholderValue string

		placeholderValue = replacement.placeholder(
			logLevel,
			replacement.placeholderValue,
		)

		rendered = strings.Replace(
			rendered, replacement.value, placeholderValue, 1,
		)
	}
	format.placeholderMutex.RUnlock()

	rendered = strings.Replace(rendered, `${prefix}`, getPrefix(prefix), 1)

	return rendered
}

func (format *Format) compile() {
	// reset compiled replacements
	format.Reset()

	var placeholder Placeholder
	var ok bool

	format.placeholderMutex.RLock()
	matches := rePlaceholder.FindAllStringSubmatch(format.formatting, -1)
	for _, match := range matches {
		var (
			replacementValue = match[0]
			placeholderName  = match[1]
			placeholderValue = match[3]
		)

		placeholder, ok = format.placeholders[placeholderName]
		if !ok {
			// placeholder with specified name not found
			continue
		}

		newReplacement := replacement{
			value:            replacementValue,
			placeholder:      placeholder,
			placeholderValue: placeholderValue,
		}

		format.replacements = append(format.replacements, newReplacement)
	}
	format.placeholderMutex.RUnlock()

	format.compiled = true
}

func getPrefix(prefix string) string {
	if prefix != "" {
		return prefix + " "
	}

	return ""
}
