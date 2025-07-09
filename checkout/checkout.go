package checkout

import (
	"errors"
	"fmt"
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
	items          map[string]int
	pricingService pricing.IPricingService
}

func NewCheckout(pricingService pricing.IPricingService) *checkout {
	return &checkout{
		items:          make(map[string]int),
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

	c.items[SKU] = c.items[SKU] + 1
	return nil
}

func (c *checkout) GetTotalPrice() (int, error) {
	if len(c.items) == 0 {
		return 0, errors.New(NoItemsError)
	}

	pricingScheme, err := c.pricingService.GetPricingScheme()
	if err != nil {
		return 0, fmt.Errorf("error occurred when getting pricing scheme. Error found %s", err.Error())
	}

	totalPrice := 0

	for sku, quantity := range c.items {
		pricingRules, ok := pricingScheme.Items[sku]
		if !ok {
			return 0, fmt.Errorf("item %s not found in pricing scheme", sku)
		}

		if quantity >= pricingRules.DiscountThreshold {
			discountGroups := quantity / pricingRules.DiscountThreshold
			itemsAfterDiscount := quantity % pricingRules.DiscountThreshold

			totalPrice += discountGroups * pricingRules.DiscountPrice
			totalPrice += itemsAfterDiscount * pricingRules.Price
		} else {
			totalPrice += quantity * pricingRules.Price
		}
	}

	return totalPrice, nil
}
