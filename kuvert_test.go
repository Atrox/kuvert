package kuvert_test

import (
	"embed"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.atrox.dev/kuvert"
	"go.atrox.dev/kuvert/themes"
)

//go:embed testdata
var testdataFS embed.FS

func TestThemeClean(t *testing.T) {
	tests := []struct {
		name  string
		email kuvert.Email

		customInstance func() *kuvert.Kuvert
	}{
		{
			name:  "empty",
			email: kuvert.Email{},
		},
		{
			name: "full",
			email: kuvert.Email{
				Summary: "This is a summary",
				Title:   "Newsletter for you!",
				Intros:  []string{"Welcome to our newsletter!", "Get ready for some exciting news."},
				Dictionary: []kuvert.Entry{
					{
						Key:   "First Name",
						Value: "John",
					},
					{
						Key:   "Last Name",
						Value: "Doe",
					},
				},
				Table: kuvert.Table{
					Data: [][]kuvert.Entry{
						{
							{
								Key:   "Product",
								Value: "Widget",
							},
							{
								Key:   "Price",
								Value: "$19.99",
							},
						},
						{
							{
								Key:   "Product",
								Value: "Gadget",
							},
							{
								Key:   "Price",
								Value: "$29.99",
							},
						},
					},
				},
				Actions: []kuvert.Action{
					{
						Instructions: "Click the button below to learn more.",
						Button: kuvert.Button{
							Link: "https://example.com",
							Text: "Learn More",
						},
					},
					{
						Instructions: "Click the button below to unsubscribe.",
						Button: kuvert.Button{
							Link: "https://example.com/unsubscribe",
							Text: "Unsubscribe",
						},
					},
				},
				Outros:    []string{"Thank you for your support!", "We hope you enjoy our products."},
				Signature: "Best regards",
			},
		},
		{
			name: "minimal",
			email: kuvert.Email{
				Intros: []string{"Verify your account"},
				Actions: []kuvert.Action{
					{
						Instructions: "Click the button below to verify your account.",
						Button: kuvert.Button{
							Link: "https://example.com/verify",
							Text: "Verify",

							Color:     "red",
							TextColor: "white",
						},
					},
				},
			},
		},
		{
			name: "table",
			email: kuvert.Email{
				Table: kuvert.Table{
					Data: [][]kuvert.Entry{
						{{Key: "Name", Value: "John"}, {Key: "Age", Value: "30"}},
						{{Key: "Name", Value: "Jane"}, {Key: "Age", Value: "25"}},
						{{Key: "Name", Value: "Alice"}, {Key: "Age", Value: "35"}},
					},
					Columns: kuvert.Columns{
						CustomWidth: map[string]string{
							"Name": "100px",
							"Age":  "50px",
						},
						CustomAlignment: map[string]string{
							"Name": "left",
							"Age":  "right",
						},
					},
				},
			},
		},
		{
			name: "summary_in_plain_text",
			email: kuvert.Email{
				Summary: "This is a summary",
				Title:   "Newsletter for you!",
				Intros:  []string{"Welcome to our newsletter!"},
				Outros:  []string{"Thank you for your support!"},
			},
			customInstance: func() *kuvert.Kuvert {
				theme, err := themes.NewCleanTheme()
				if err != nil {
					t.Fatal(err)
				}
				theme.SummaryInPlainText = true
				return kuvert.New(theme, kuvert.Branding{
					Name: "Test Product",
					Link: "https://example.com",
				})
			},
		},
	}

	instance := kuvertInstance(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := instance
			if tt.customInstance != nil {
				instance = tt.customInstance()
			}

			output, err := instance.Render(tt.email)
			if err != nil {
				t.Fatal(err)
			}

			local := getLocalRenderedEmail(t, "clean", tt.name)

			if !assert.Equal(t, local.HTML, output.HTML) {
				_ = os.WriteFile(fmt.Sprintf("./testdata/clean/%s.html", tt.name), []byte(output.HTML), 0o644)
			}
			if !assert.Equal(t, local.PlainText, output.PlainText) {
				_ = os.WriteFile(fmt.Sprintf("./testdata/clean/%s.txt", tt.name), []byte(output.PlainText), 0o644)
			}
		})
	}
}

func BenchmarkCleanTheme(b *testing.B) {
	instance := kuvertInstance(b)
	email := kuvert.Email{
		Summary: "This is a summary",
		Title:   "Newsletter for you!",
		Intros:  []string{"Welcome to our newsletter!", "Get ready for some exciting news."},
		Dictionary: []kuvert.Entry{
			{
				Key:   "First Name",
				Value: "John",
			},
			{
				Key:   "Last Name",
				Value: "Doe",
			},
		},
		Table: kuvert.Table{
			Data: [][]kuvert.Entry{
				{
					{
						Key:   "Product",
						Value: "Widget",
					},
					{
						Key:   "Price",
						Value: "$19.99",
					},
				},
				{
					{
						Key:   "Product",
						Value: "Gadget",
					},
					{
						Key:   "Price",
						Value: "$29.99",
					},
				},
			},
		},
		Actions: []kuvert.Action{
			{
				Instructions: "Click the button below to learn more.",
				Button: kuvert.Button{
					Link: "https://example.com",
					Text: "Learn More",
				},
			},
			{
				Instructions: "Click the button below to unsubscribe.",
				Button: kuvert.Button{
					Link: "https://example.com/unsubscribe",
					Text: "Unsubscribe",
				},
			},
		},
		Outros:    []string{"Thank you for your support!", "We hope you enjoy our products."},
		Signature: "Best regards",
	}

	b.Run("Render", func(b *testing.B) {
		for range b.N {
			_, err := instance.Render(email)
			if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("RenderSequential", func(b *testing.B) {
		for range b.N {
			_, err := instance.RenderHTML(email)
			if err != nil {
				b.Error(err)
			}
			_, err = instance.RenderPlainText(email)
			if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("RenderHTML", func(b *testing.B) {
		for range b.N {
			_, err := instance.RenderHTML(email)
			if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("RenderPlainText", func(b *testing.B) {
		for range b.N {
			_, err := instance.RenderPlainText(email)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func ExampleNew() {
	theme, err := themes.NewDefaultTheme()
	if err != nil {
		log.Fatal(err)
	}

	branding := kuvert.Branding{
		Name: "Test Product",
		Link: "https://example.com",
	}

	instance := kuvert.New(theme, branding)

	email := kuvert.Email{
		Title: "Welcome to Test Product",
		Intros: []string{
			"We're excited to have you on board.",
		},
		Actions: []kuvert.Action{
			{
				Instructions: "Click the button below to verify your email address.",
				Button: kuvert.Button{
					Link: "https://example.com/verify",
					Text: "Verify Email",
				},
			},
		},
		Outros: []string{
			"Need help, or have questions? Just reply to this email, we'd love to help.",
		},
	}

	renderedEmail, err := instance.Render(email)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(renderedEmail.HTML)
	fmt.Println(renderedEmail.PlainText)
}

func kuvertInstance(t testing.TB) *kuvert.Kuvert {
	theme, err := themes.NewCleanTheme()
	if err != nil {
		t.Fatal(err)
	}
	branding := kuvert.Branding{
		Name: "Test Product",
		Link: "https://example.com",
	}
	return kuvert.New(theme, branding)
}

func getLocalRenderedEmail(t testing.TB, theme, name string) *kuvert.RenderedEmail {
	htmlPath := fmt.Sprintf("testdata/%s/%s.html", theme, name)
	plainTextPath := fmt.Sprintf("testdata/%s/%s.txt", theme, name)

	html, err := testdataFS.ReadFile(htmlPath)
	if err != nil {
		t.Error(err)
	}
	plainText, err := testdataFS.ReadFile(plainTextPath)
	if err != nil {
		t.Error(err)
	}

	return &kuvert.RenderedEmail{
		HTML:      string(html),
		PlainText: string(plainText),
	}
}
