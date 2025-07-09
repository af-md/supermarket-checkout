package checkout

import "errors"

const (
	EmptySKUError      = "SKU cannot be empty"
	WhitespaceSKUError = "SKU cannot have whitespace"
)

type ICheckout interface {
	Scan(SKU string) (err error)
	GetTotalPrice() (totalPrice int, err error)
}

type checkout struct {
	items []string
}

type item struct {
	SKU string
}

func NewCheckout() *checkout {
	return &checkout{
		items: make([]string, 0),
	}
}

func (c *checkout) Scan(SKU string) error {
	// could have trimmed the whitespace here, but I wanted to show how I would handle multiple edge cases
	if SKU == "" {
		return errors.New(EmptySKUError)
	}

	if SKU == " " {
		return errors.New(WhitespaceSKUError)
	}

	c.items = append(c.items, SKU)
	return nil
}

func (c *checkout) GetTotalPrice() (int, error) {
	panic("not implemented")
}
