package checkout

import (
	"errors"
	"supermarket-checkout/pricing"
)

const (
	// SKU related error messages
	EmptySKUError      = "SKU cannot be empty"
	WhitespaceSKUError = "SKU cannot have whitespace"

	// Generate price related error messages
	NoItemsError = "no items scanned, checkout empty"
)

type ICheckout interface {
	Scan(SKU string) (err error)
	GetTotalPrice() (totalPrice int, err error)
}

type checkout struct {
	items          []string
	pricingService pricing.IPricingService
}

func NewCheckout(pricingService pricing.IPricingService) *checkout {
	return &checkout{
		items:          make([]string, 0),
		pricingService: pricingService,
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
	if len(c.items) == 0 {
		return 0, errors.New(NoItemsError)
	}

	totalPrice, err := c.pricingService.GetPrice(c.items[0])
	if err != nil {
		return 0, err
	}

	return totalPrice, nil
}
