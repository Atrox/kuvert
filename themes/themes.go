package themes

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
	return string(FileDefaultTemplateGohtml)
}

func (*defaultTheme) PlainTextTemplate() string {
	return string(FileDefaultTemplateGotext)
}
