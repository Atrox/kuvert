package themes

import _ "embed"

//go:embed default/template.gohtml
var defaultThemeHtml string

//go:embed default/template.gotext
var defaultThemeText string

// Theme is an interface to implement when creating a new theme
type Theme interface {
	Name() string              // The name of the theme
	HTMLTemplate() string      // The golang template for HTML emails
	PlainTextTemplate() string // The golang template for plain text emails (can be basic HTML)
}

func Default() *defaultTheme {
	return &defaultTheme{}
}

// Default is the theme by default
type defaultTheme struct{}

func (*defaultTheme) Name() string {
	return "default"
}

func (*defaultTheme) HTMLTemplate() string {
	return defaultThemeHtml
}

func (*defaultTheme) PlainTextTemplate() string {
	return defaultThemeText
}
