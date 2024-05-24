package themes

import (
	"go.atrox.dev/kuvert/themes/clean"
)

// NewDefaultTheme creates a new default theme for kuvert.
// Default theme is the [clean.Theme] theme.
func NewDefaultTheme() (*clean.Theme, error) {
	return NewCleanTheme()
}

// NewCleanTheme creates a new clean theme for kuvert.
func NewCleanTheme() (*clean.Theme, error) {
	return clean.New()
}
