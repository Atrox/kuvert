// Package kuvert is a simple email generator for Go.
// It is designed to be easy to use and to allow you to quickly generate beautiful emails.
//
// Kuvert is designed to be used with a [Theme].
// A theme is responsible for rendering the email in both HTML and plain text.
// You can create your own theme or use one of the built-in themes.
//
// Kuvert supports your company branding.
// [Branding] is a struct that contains information about your company product (brand).
// It appears in the header and footer of the emails.
package kuvert

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

// Kuvert is an instance of the kuvert email generator.
type Kuvert struct {
	Theme    Theme
	Branding *Branding
}

// Branding represents your company product (brand).
// Appears in header & footer of e-mails.
type Branding struct {
	Name        string
	Link        string
	Logo        string
	Copyright   string
	TroubleText string
}

// New creates a new [Kuvert] instance.
func New(theme Theme, branding Branding) *Kuvert {
	if branding.Copyright == "" {
		branding.Copyright = fmt.Sprintf(
			"Copyright © %s %s. All rights reserved.",
			strconv.Itoa(time.Now().Year()), branding.Name,
		)
	}
	if branding.TroubleText == "" {
		branding.TroubleText = "If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser."
	}

	return &Kuvert{
		Theme:    theme,
		Branding: &branding,
	}
}

// RenderedEmail contains the rendered HTML and plain text of an email.
type RenderedEmail struct {
	HTML      string
	PlainText string
}

// Render renders the email in both HTML and plain text format.
func (k *Kuvert) Render(email Email) (*RenderedEmail, error) {
	renderedEmail := &RenderedEmail{}
	wg := &errgroup.Group{}
	wg.Go(func() error {
		var err error
		renderedEmail.HTML, err = k.RenderHTML(email)
		return err
	})
	wg.Go(func() error {
		var err error
		renderedEmail.PlainText, err = k.RenderPlainText(email)
		return err
	})
	if err := wg.Wait(); err != nil {
		return nil, err
	}

	return renderedEmail, nil
}

// RenderHTML renders the email in HTML format.
func (k *Kuvert) RenderHTML(email Email) (string, error) {
	email.fillDefaults()

	return k.Theme.HTML(k, &email)
}

// RenderPlainText renders the email in plain text format.
func (k *Kuvert) RenderPlainText(email Email) (string, error) {
	email.fillDefaults()

	return k.Theme.PlainText(k, &email)
}
