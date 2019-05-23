//go:generate fileb0x filebox.json
package kuvert

import (
	"fmt"

	"go.atrox.dev/kuvert/themes"
)

var instance *Kuvert

// Kuvert is an instance of the kuvert email generator
type Kuvert struct {
	Theme   themes.Theme
	Product *Product
}

// Product represents your company product (brand)
// Appears in header & footer of e-mails
type Product struct {
	Name        string
	Link        string
	Logo        string
	Copyright   string
	TroubleText string
}

func New(product *Product) *Kuvert {
	if product.Copyright == "" {
		product.Copyright = fmt.Sprintf("Copyright © 2019 %s. All rights reserved.", product.Name)
	}
	if product.TroubleText == "" {
		product.TroubleText = "If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser."
	}

	return &Kuvert{
		Theme:   themes.Default(),
		Product: product,
	}
}

func (k *Kuvert) NewEmail() *Email {
	return &Email{
		Kuvert: k,

		Intros:     []string{},
		Dictionary: []Entry{},
		Outros:     []string{},
		Signature:  "Yours truly",
		Greeting:   "Hi",
	}
}

func Instance() *Kuvert {
	return instance
}

func Init(product *Product) {
	if product.Copyright == "" {
		product.Copyright = fmt.Sprintf("Copyright © 2019 %s. All rights reserved.", product.Name)
	}
	if product.TroubleText == "" {
		product.TroubleText = "If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser."
	}

	instance.Product = product
}

func SetProduct(product *Product) {
	instance.Product = product
}

func SetTheme(theme themes.Theme) {
	instance.Theme = theme
}

func NewEmail() *Email {
	if instance.Theme == nil {
		instance.Theme = themes.Default()
	}
	if instance.Product == nil {
		instance.Product = &Product{}
	}

	return instance.NewEmail()
}
