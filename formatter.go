package lorg

type Formatter interface {
	SetPlaceholders(placeholders map[string]Placeholder)
	SetPlaceholder(name string, placeholder Placeholder)
	GetPlaceholders() map[string]Placeholder
	Render(logLevel Level) string
	Reset()
}
