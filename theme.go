package kuvert

// Theme is the interface that all themes must implement
// Themes are responsible for rendering the email in both HTML and plain text.
type Theme interface {
	HTML(kuvert *Kuvert, email *Email) (string, error)
	PlainText(kuvert *Kuvert, email *Email) (string, error)
}
