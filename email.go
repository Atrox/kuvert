package kuvert

import (
	"bytes"
	"html/template"

	"github.com/Masterminds/sprig"
	"github.com/jaytaylor/html2text"
)

var templateFuncs = template.FuncMap{
	"url": func(s string) template.URL {
		return template.URL(s)
	},
	"html": func(s string) template.HTML {
		return template.HTML(s)
	},
}

// Email is the email containing a body
type Email struct {
	Name       string   // The name of the contacted person
	Intros     []string // Intro sentences, first displayed in the email
	Dictionary []Entry  // A list of key+value (useful for displaying parameters/settings/personal info)
	Table      Table    // Table is an table where you can put data (pricing grid, a bill, and so on)
	Actions    []Action // Actions are a list of actions that the user will be able to execute via a button click
	Outros     []string // Outro sentences, last displayed in the email
	Greeting   string   // Greeting for the contacted person (default to 'Hi')
	Signature  string   // Signature for the contacted person (default to 'Yours truly')
	Title      string   // Title replaces the greeting+name when set

	Kuvert *Kuvert
}

// Entry is a simple entry of a map
// Allows using a slice of entries instead of a map
// Because Golang maps are not ordered
type Entry struct {
	Key   string
	Value string
}

// Table is an table where you can put data (pricing grid, a bill, and so on)
type Table struct {
	Data    [][]Entry // Contains data
	Columns Columns   // Contains meta-data for display purpose (width, alignement)
}

// Columns contains meta-data for the different columns
type Columns struct {
	CustomWidth     map[string]string
	CustomAlignment map[string]string
}

// Action is an action the user can do on the email (click on a button)
type Action struct {
	Instructions string
	Button       Button
}

// Button defines an action to launch
type Button struct {
	Color     string
	TextColor string
	Text      string
	Link      string
}

// GenerateHTML generates the email body from data to an HTML Reader
// This is for modern email clients
func (e *Email) GenerateHTML() (string, error) {
	return e.generateTemplate(e.Kuvert.Theme.HTMLTemplate())
}

// GeneratePlainText generates the email body from data
// This is for old email clients
func (e *Email) GeneratePlainText() (string, error) {
	tmpl, err := e.generateTemplate(e.Kuvert.Theme.PlainTextTemplate())
	if err != nil {
		return "", err
	}
	return html2text.FromString(tmpl, html2text.Options{PrettyTables: true})
}

func (e *Email) generateTemplate(tmpl string) (string, error) {
	// Generate the email from Golang template
	// Allow usage of simple function from sprig : https://github.com/Masterminds/sprig
	t, err := template.New("kuvert").Funcs(sprig.FuncMap()).Funcs(templateFuncs).Parse(tmpl)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = t.Execute(&b, e)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
