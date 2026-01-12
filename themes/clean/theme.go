package clean

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"

	"github.com/inbucket/html2text"
	"github.com/vanng822/go-premailer/premailer"
	"go.atrox.dev/kuvert"
)

//go:embed template.gohtml
var htmlTemplate string

//go:embed template.gotext
var plainTextTemplate string

type templateContext struct {
	Branding *kuvert.Branding
	Email    *kuvert.Email
	Theme    *Theme
}

type Theme struct {
	SummaryInPlainText bool // SummaryInPlainText determines if the summary should be included in the plain text version of the email

	htmlTemplate      *template.Template
	plainTextTemplate *template.Template
}

func New() (*Theme, error) {
	templateFuncs := template.FuncMap{
		"replace": func(old, new, src string) string {
			return strings.ReplaceAll(src, old, new)
		},
	}

	htmlTemplate, err := template.New("html").Funcs(templateFuncs).Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}

	plainTextTemplate, err := template.New("plainText").Funcs(templateFuncs).Parse(plainTextTemplate)
	if err != nil {
		return nil, err
	}

	return &Theme{
		htmlTemplate:      htmlTemplate,
		plainTextTemplate: plainTextTemplate,
	}, nil
}

func (t *Theme) HTML(kuvert *kuvert.Kuvert, email *kuvert.Email) (string, error) {
	buf := &bytes.Buffer{}
	err := t.htmlTemplate.Execute(buf, &templateContext{
		Branding: kuvert.Branding,
		Email:    email,
		Theme:    t,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}
	tmpl := buf.String()

	prem, err := premailer.NewPremailerFromString(tmpl, premailer.NewOptions())
	if err != nil {
		return "", fmt.Errorf("failed to create premailer: %w", err)
	}
	return prem.Transform()
}

func (t *Theme) PlainText(kuvert *kuvert.Kuvert, email *kuvert.Email) (string, error) {
	buf := &bytes.Buffer{}
	err := t.plainTextTemplate.Execute(buf, &templateContext{
		Branding: kuvert.Branding,
		Email:    email,
		Theme:    t,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute plain text template: %w", err)
	}

	return html2text.FromReader(buf, html2text.Options{PrettyTables: true})
}
