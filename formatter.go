package lorg

// Formatter is the interface which implemented by Format structure, it's
// usable if you want to create your own formating mechanism.
type Formatter interface {

	// SetPlaceholders sets specified placeholders in given formatter.
	// map key is a name of placeholder.
	SetPlaceholders(placeholders map[string]Placeholder)

	// SetPlaceholder sets specified placeholder with specified name in given
	// formatter
	SetPlaceholder(name string, placeholder Placeholder)

	// GetPlaceholders gets placeholders of given formatter.
	GetPlaceholders() map[string]Placeholder

	// Log.Print/Log.Warning/Log.Error/etc. will calls formatter Render
	// Render should return string for Log instance which will be used for
	// fmt.Sprintf of logging record.
	Render(logLevel Level, prefix string) string

	// Reset Formatter state.
	Reset()
}
